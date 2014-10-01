package tunnel

import (
	"fmt"
	"io"
	"net"

	"encrypt"

	"github.com/hamo/golog"
)

type RawSocketClient struct {
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

type RawSocketServer struct {
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

func NewRawSocketClient(addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*RawSocketClient, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &RawSocketClient{
		addr:      addr,
		port:      port,
		crypto:    c,
		ivReady:   make(chan bool, 0),
		masterKey: c.GenKey(masterKey),
		password:  c.GenKey(password),
		logger:    logger,
	}, nil
}

func NewRawSocketServer(masterKey string, encryptMethod string, logger *golog.GoLogger) (*RawSocketServer, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &RawSocketServer{
		crypto:    c,
		ivReady:   make(chan bool, 0),
		masterKey: c.GenKey(masterKey),
		logger:    logger,
	}, nil
}

func (r *RawSocketClient) Dial() error {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", r.addr, r.port))
	r.c = c
	return err
}

func (r *RawSocketServer) Accept(c net.Conn) {
	r.c = c
}

func (r *RawSocketClient) Close() {
	r.c.Close()
}

func (r *RawSocketServer) Close() {
	r.c.Close()
}

//FIXME
func (r *RawSocketServer) SetPassword(passwd string) {
	r.password = r.crypto.GenKey(passwd)
}

func (r *RawSocketClient) ReadMaster(buf []byte, full bool) (int, error) {
	panic("client call readMaster")
}

func (r *RawSocketServer) ReadMaster(buf []byte, full bool) (int, error) {
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

func (r *RawSocketClient) ReadUser(buf []byte, full bool) (int, error) {
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

func (r *RawSocketServer) ReadUser(buf []byte, full bool) (int, error) {
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

func (r *RawSocketClient) ReadContent(buf []byte) (int, error) {
	if r.userEncryptR == nil {
		// call readuser from client, wait for iv ready
		<-r.ivReady
		r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := make([]byte, len(buf))
	var n int
	var err error

	n, err = r.c.Read(ciphertext)

	plaintext := r.userEncryptR.Decrypt(ciphertext[:n])

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *RawSocketServer) ReadContent(buf []byte) (int, error) {
	// call after ParseUserHeader
	if r.userEncryptR == nil {
		r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := make([]byte, len(buf))
	var n int
	var err error

	n, err = r.c.Read(ciphertext)

	plaintext := r.userEncryptR.Decrypt(ciphertext[:n])

	copy(buf[:n], plaintext[:n])
	return n, err
}

func (r *RawSocketClient) WriteMaster(buf []byte) (int, error) {
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

func (r *RawSocketServer) WriteMaster(buf []byte) (int, error) {
	panic("Server call WriteMaster")
}

func (r *RawSocketClient) WriteUser(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}

func (r *RawSocketServer) WriteUser(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		<-r.ivReady
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}

func (r *RawSocketClient) WriteContent(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}

func (r *RawSocketServer) WriteContent(buf []byte) (int, error) {
	if r.userEncryptW == nil {
		<-r.ivReady
		r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}
