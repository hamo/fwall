package tunnel

import (
	"fmt"
	"github.com/hamo/fwall/encrypt"
	"github.com/hamo/golog"
	"io"
	"net"
)

type Reader interface {
	// Read IV and Master header
	ReadMaster(p []byte, full bool) (int, error)

	// Read User header
	ReadUser(p []byte, full bool) (int, error)

	// Read Content data
	ReadContent(p []byte) (int, error)
}

type Writer interface {
	// Write IV and Master header
	WriteMaster(p []byte) (n int, err error)

	// Write User header
	WriteUser(p []byte) (n int, err error)

	// Write Content data
	WriteContent(p []byte) (n int, err error)
}

type ProxyAgent interface {
	Dial() error
	Accept(c net.Conn)
	Close()
	SetPassword(passwd string)

	// Read IV and Master header
	ReadMaster(p []byte, full bool) (int, error)

	// Read User header
	ReadUser(p []byte, full bool) (int, error)

	// Read Content data
	ReadContent(p []byte) (int, error)

	// Write IV and Master header
	WriteMaster(p []byte) (n int, err error)

	// Write User header
	WriteUser(p []byte) (n int, err error)

	// Write Content data
	WriteContent(p []byte) (n int, err error)
}

func NewClient(clientType string, addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (ProxyAgent, error) {
	switch clientType {
	case "lz4":
		return NewLZ4SocketClient(addr, port, masterKey, encryptMethod, password, logger)
	case "plain":
		return NewPlainTunnelClient(addr, port, masterKey, encryptMethod, password, logger)
	default:
		return NewRawSocketClient(addr, port, masterKey, encryptMethod, password, logger)
	}
}

func NewServer(serverType string, masterKey string, encryptMethod string, logger *golog.GoLogger) (ProxyAgent, error) {
	switch serverType {
	case "lz4":
		return NewLZ4SocketServer(masterKey, encryptMethod, logger)
	case "plain":
		return NewPlainTunnelServer(masterKey, encryptMethod, logger)
	default:
		return NewRawSocketServer(masterKey, encryptMethod, logger)
	}
}

// Embed Base behavior in your tunnel. Overriding them when needed.
// Hope the following job will reduce the duplicate codes.
type ClientBase struct {
	addr string
	port int

	iv      []byte
	ivReady chan bool

	masterKey []byte
	password  []byte

	crypto        *encrypt.CryptoInfo
	masterEncrypt encrypt.Encryption
	userEncryptR  encrypt.Encryption
	userEncryptW  encrypt.Encryption

	c net.Conn

	logger *golog.GoLogger
}

type ServerBase struct {
	iv      []byte
	ivReady chan bool

	masterKey []byte
	password  []byte

	crypto        *encrypt.CryptoInfo
	masterEncrypt encrypt.Encryption
	userEncryptR  encrypt.Encryption
	userEncryptW  encrypt.Encryption

	c net.Conn

	logger *golog.GoLogger
}

func (r *ClientBase) Dial() error {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", r.addr, r.port))
	c.(*net.TCPConn).SetNoDelay(false)
	r.c = c
	return err
}

func (r *ServerBase) Accept(c net.Conn) {
	r.c = c
}

func (r *ClientBase) Close() {
	r.c.Close()
}

func (r *ServerBase) Close() {
	r.c.Close()
}

func (r *ServerBase) Dial() error {
	panic("server call Dial")
}

func (r *ClientBase) Accept(c net.Conn) {
	panic("client call Accept")
}

func (r *ClientBase) SetPassword(passwd string) {
	panic("client call SetPassword")
}

func (r *ServerBase) SetPassword(passwd string) {
	r.password = r.crypto.GenKey(passwd)
}

func (r *ClientBase) ReadMaster(buf []byte, full bool) (int, error) {
	panic("client call readMaster")
}

func (r *ServerBase) ReadMaster(buf []byte, full bool) (int, error) {
	if r.masterEncrypt == nil {
		// First call readMaster, so get IV first
		iv := make([]byte, r.crypto.IVlen)
		n, err := io.ReadFull(r.c, iv)
		if err != nil {
			close(r.ivReady)
			return n, err
		}
		r.logger.Debugf("get IV: %v, err : %v", iv, err)
		r.iv = iv
		close(r.ivReady)
		r.masterEncrypt = r.crypto.NewCrypto(r.iv, r.masterKey)
	}

	ciphertext := make([]byte, len(buf))
	var n int
	var err error

	if full {
		n, err = io.ReadFull(r.c, ciphertext)
	} else {
		n, err = r.c.Read(ciphertext)
	}

	plaintext := r.masterEncrypt.Decrypt(ciphertext[:n])

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *ClientBase) ReadUser(buf []byte, full bool) (int, error) {
	if r.userEncryptR == nil {
		// call readuser from client, wait for iv ready
		<-r.ivReady
		r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := make([]byte, len(buf))
	var n int
	var err error

	if full {
		n, err = io.ReadFull(r.c, ciphertext)
	} else {
		n, err = r.c.Read(ciphertext)
	}

	plaintext := r.userEncryptR.Decrypt(ciphertext[:n])

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *ServerBase) ReadUser(buf []byte, full bool) (int, error) {
	// call after readMaster
	if r.userEncryptR == nil {
		r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := make([]byte, len(buf))
	var n int
	var err error

	if full {
		n, err = io.ReadFull(r.c, ciphertext)
	} else {
		n, err = r.c.Read(ciphertext)
	}

	plaintext := r.userEncryptR.Decrypt(ciphertext[:n])

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *ClientBase) ReadContent(buf []byte) (int, error) {
	plaintext := make([]byte, len(buf))
	var n int
	var err error

	n, err = r.c.Read(plaintext)
	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *ServerBase) ReadContent(buf []byte) (int, error) {
	// call after ParseUserHeader
	plaintext := make([]byte, len(buf))
	var n int
	var err error

	n, err = r.c.Read(plaintext)

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *ClientBase) WriteMaster(buf []byte) (int, error) {
	if r.masterEncrypt == nil {
		// first time to write master
		r.iv = r.crypto.GenIV()
		close(r.ivReady)
		r.logger.Debugf("first write. IV: %v", r.iv)
		r.c.Write(r.iv)

		r.masterEncrypt = r.crypto.NewCrypto(r.iv, r.masterKey)
	}

	return r.c.Write(r.masterEncrypt.Encrypt(buf))
}

func (r *ServerBase) WriteMaster(buf []byte) (int, error) {
	panic("Server call WriteMaster")
}

func (r *ClientBase) WriteUser(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}

func (r *ServerBase) WriteUser(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		<-r.ivReady
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}

func (r *ClientBase) WriteContent(buf []byte) (int, error) {
	return r.c.Write(buf)
}

func (r *ServerBase) WriteContent(buf []byte) (int, error) {
	return r.c.Write(buf)
}
