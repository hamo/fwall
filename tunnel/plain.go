package tunnel

import (
	"fmt"

	"encrypt"

	"github.com/hamo/golog"
)

type PlainTunnelClient struct {
	ClientBase
}

type PlainTunnelServer struct {
	ServerBase
}

func NewPlainTunnelClient(addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*RawSocketClient, error) {
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

func NewPlainTunnelServer(masterKey string, encryptMethod string, logger *golog.GoLogger) (*RawSocketServer, error) {
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
