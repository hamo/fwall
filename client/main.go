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
	defer c.Close()
	if err != nil {
		logger.Debugf("handShake err: %s", err)
		return
	}

	commandCode, addressType, address, port, err := parseReq(c)
	if err != nil {
		logger.Warningf("parseReq failed: %s", err)
		return
	}

	logger.Debugf("commandCode: %d\n", commandCode)
	logger.Debugf("addressType: %d\n", addressType)
	logger.Debugf("port: %d\n", port)

	_, err = c.Write(reqAnswer)
	if err != nil {
		logger.Warningf("write req answer failed: %s", err)
		return
	}

	proxyAgent, err := tunnel.NewClient(lc.Tunnel, lc.Server, lc.ServerPort, lc.MasterKey, lc.EncryptMethod, lc.Password, logger)
	if err != nil {
		logger.Warningf("Create tunnel failed: %s", err)
		return
	}

	err = proxyAgent.Dial()
	if err != nil {
		logger.Warningf("Dial to %s:%d failed: %s", lc.Server, lc.ServerPort, err)
		return
	}
	defer proxyAgent.Close()

	client := protocol.NewClient(lc.Username, addressType, address, port, logger)
	go client.Upstream(c, proxyAgent)
	client.Downstream(c, proxyAgent)
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
