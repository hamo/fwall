package tunnel

import (
	"fmt"

	"github.com/hamo/fwall/encrypt"

	"github.com/hamo/golog"
)

type PlainTunnelClient struct {
	ClientBase
}

type PlainTunnelServer struct {
	ServerBase
}

func NewPlainTunnelClient(addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*PlainTunnelClient, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &PlainTunnelClient{
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

func NewPlainTunnelServer(masterKey string, encryptMethod string, logger *golog.GoLogger) (*PlainTunnelServer, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &PlainTunnelServer{
		ServerBase{
			crypto:    c,
			ivReady:   make(chan bool, 0),
			masterKey: c.GenKey(masterKey),
			logger:    logger,
		},
	}, nil
}
