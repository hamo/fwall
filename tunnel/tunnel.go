package tunnel

import (
	"net"
	"github.com/hamo/golog"
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
	Dial() (error)
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
	default:
		return NewRawSocketClient(addr, port, masterKey, encryptMethod, password, logger)
	}
}

func NewServer(serverType string, masterKey string, encryptMethod string, logger *golog.GoLogger) (ProxyAgent, error) {
	switch serverType {
	default:
		return NewRawSocketServer(masterKey, encryptMethod, logger)
	}
}
