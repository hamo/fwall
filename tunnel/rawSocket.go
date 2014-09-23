package tunnel

import (
	"fmt"
	"io"
	"net"
	"runtime"

	"encrypt"

	"github.com/hamo/golog"
)

type RawSocket struct {
	addr string
	port int

	side string

	iv        []byte
	masterKey []byte
	password  []byte

	crypto        *encrypt.CryptoInfo
	masterEncrypt encrypt.Encryption
	userEncryptR  encrypt.Encryption
	userEncryptW  encrypt.Encryption

	c net.Conn

	logger *golog.GoLogger
}

func NewRawSocket(addr string, port int, side string, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*RawSocket, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &RawSocket{
		addr:      addr,
		port:      port,
		side:      side,
		crypto:    c,
		masterKey: c.GenKey(masterKey),
		password:  c.GenKey(password),
		logger:    logger,
	}, nil
}

func (r *RawSocket) Dial() error {
	c, err := net.Dial("tcp", fmt.Sprintf("%s:%d", r.addr, r.port))
	r.c = c
	return err
}

func (r *RawSocket) Accept(c net.Conn) {
	r.c = c
}

func (r *RawSocket) Close() {
	r.c.Close()
}

//FIXME
func (r *RawSocket) SetPassword(passwd string) {
	r.password = r.crypto.GenKey(passwd)
}

func (r *RawSocket) ReadMaster(buf []byte, full bool) (int, error) {
	if r.side == "client" {
		panic("client call readMaster")
	}
	if r.masterEncrypt == nil {
		// First call readMaster, so get IV first
		iv := make([]byte, r.crypto.IVlen)
		// FIXME: check err
		_, err := io.ReadFull(r.c, iv)
		r.logger.Debugf("get IV: %v, err : %v", iv, err)
		r.iv = iv
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

	copy(buf, plaintext)
	return n, err
}

func (r *RawSocket) ReadUser(buf []byte, full bool) (int, error) {
	switch r.side {
	case "client":
		if r.userEncryptR == nil {
			// call readuser from client, wait for iv ready
			for true {
				if len(r.iv) == r.crypto.IVlen {
					break
				}
				runtime.Gosched()
			}
			r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
		}
	case "server":
		// call after readMaster
		if r.userEncryptR == nil {
			r.userEncryptR = r.crypto.NewCrypto(r.iv, r.password)
		}
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

	copy(buf, plaintext)
	return n, err
}

func (r *RawSocket) WriteMaster(buf []byte) (int, error) {
	if r.side == "server" {
		panic("call write master from server side")
	}
	if r.masterEncrypt == nil {
		// first time to write master
		r.iv = r.crypto.GenIV()
		r.logger.Debugf("first write. IV: %v", r.iv)
		r.c.Write(r.iv)

		r.masterEncrypt = r.crypto.NewCrypto(r.iv, r.masterKey)
	}

	return r.c.Write(r.masterEncrypt.Encrypt(buf))
}

func (r *RawSocket) WriteUser(buf []byte) (int, error) {
	switch r.side {
	case "client":
		if r.userEncryptW == nil {
			r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
		}
	case "server":
		if r.userEncryptW == nil {
			for true {
				if len(r.iv) == r.crypto.IVlen {
					break
				}
				runtime.Gosched()
			}
			r.userEncryptW = r.crypto.NewCrypto(r.iv, r.password)
		}
	}

	ciphertext := r.userEncryptW.Encrypt(buf)

	return r.c.Write(ciphertext)
}
