package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"protocol"
	"tunnel"

	"github.com/hamo/golog"
)

const (
	version = "0.1"
)

var (
	lc *localConfig
)

var (
	flDebug      *bool
	flConfigFile *string

	logger *golog.GoLogger
)

func init() {
	flDebug = flag.Bool("d", false, "debug switch")
	flConfigFile = flag.String("c", "./config.json", "config file")
}

func handleTCPConnection(c net.Conn) {
	err := handShake(c)

	if err != nil {
		logger.Debugf("handShake err: %s", err)
		c.Close()
		return
	}

	commandCode, addressType, address, port, err := parseReq(c)
	if err != nil {
		logger.Fatalf("parseReq failed: %s", err)
	}

	logger.Debugf("commandCode: %d\n", commandCode)
	logger.Debugf("addressType: %d\n", addressType)
	logger.Debugf("port: %d\n", port)

	c.Write(reqAnswer)

	// FIXME: switch will led a compile error.
	// FIXME2: if we declare proxyAgent before switch the call we must set a
	// unmodifiable type. That's why I hate strong typed language!!!
	// However if we create a new interface (no...), it must include as much
	// functions as possible. E.g.: Close(), all the read/write functions.

	var proxyAgent tunnel.Agent
	switch lc.Tunnel {
	case "Raw":
		proxyAgent, err = tunnel.RawSocketDial(lc.Server, lc.ServerPort)
	case "Http":
		proxyAgent, err = tunnel.HttpTunnelDial(lc.Server, lc.ServerPort)
	default:
		logger.Warningf("Unsupported tunnel type: %s\nUsing Raw as default", lc.Tunnel)
		proxyAgent, err = tunnel.RawSocketDial(lc.Server, lc.ServerPort)
	}

	if err != nil {
		logger.Warningf("Dial to %s:%d failed: %s", lc.Server, lc.ServerPort, err)
		return
	}
	defer proxyAgent.Close()

	client := protocol.NewClient(lc.Username, addressType, address, port, logger)
	go client.Upstream(c, proxyAgent)
	client.Downstream(c, proxyAgent)

	proxyAgent.Close()
	c.Close()
}

func main() {
	var err error

	// FIXME: configurable logger file
	logger = golog.New(os.Stdout)

	flag.Parse()

	logger.SetDebug(*flDebug)

	lc, err = parseConfigFile(*flConfigFile)
	if err != nil {
		logger.Fatalf("Parse config file err: %s", err)
	}

	logger.Infof("fwall started. Version: %s", version)

	lnTCP, err := net.Listen("tcp", fmt.Sprintf(":%d", lc.LocalPort))
	if err != nil {
		logger.Fatalf("Listen to socks5 port failed: %s", err)
	}
	defer lnTCP.Close()

	for {
		connTCP, err := lnTCP.Accept()
		if err != nil {
			logger.Debugf("Accept return err: %s", err)
			continue
		}

		go handleTCPConnection(connTCP)
	}

}
