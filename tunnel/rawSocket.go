package tunnel

import (
	"fmt"

	"encrypt"

	"github.com/hamo/golog"
)

type RawSocketClient struct {
	ClientBase
}

type RawSocketServer struct {
	ServerBase
}

func NewRawSocketClient(addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*RawSocketClient, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &RawSocketClient{
		ClientBase{
			addr:      addr,
			port:      port,
			crypto:    c,
			ivReady:   make(chan bool, 0),
			masterKey: c.GenKey(masterKey),
			password:  c.GenKey(password),
			logger:    logger,
		},
	}, nil
}

func NewRawSocketServer(masterKey string, encryptMethod string, logger *golog.GoLogger) (*RawSocketServer, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &RawSocketServer{
		ServerBase{
			crypto:    c,
			ivReady:   make(chan bool, 0),
			masterKey: c.GenKey(masterKey),
			logger:    logger,
		},
	}, nil
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
