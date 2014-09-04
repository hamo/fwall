package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

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

func fmtHeader(remoteAddr []byte) []byte {
	lengthByte := []byte{byte(len(remoteAddr))}
	fmtedHeader := append(lengthByte, remoteAddr...)
	return fmtedHeader
}

func handleTCPConnection(c net.Conn) {
	err := handShake(c)

	if err != nil {
		logger.Debugf(*flDebug, "handShake err: %s", err)
		c.Close()
		return
	}

	commandCode, addressType, address, port, err := parseReq(c)
	if err != nil {
		logger.Fatalf("parseReq failed: %s", err)
	}

	logger.Debugf(*flDebug, "commandCode: %d\n", commandCode)
	logger.Debugf(*flDebug, "addressType: %d\n", addressType)
	logger.Debugf(*flDebug, "port: %d\n", port)

	c.Write(reqAnswer)

	// Move to server
	//	realAddr := string(address)
	var realAddr string
	if addressType == 0x03 {
		realAddr = string((address.Bytes()[1:]))
	}

	realAddr = realAddr + ":" + strconv.Itoa(int(port))
	remoteAddr := fmt.Sprintf("%s:%d", lc.Server, lc.ServerPort)
	fmtedHeader := fmtHeader([]byte(realAddr))

	proxyAgent, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		logger.Warningf("Dial to %s failed: %s", remoteAddr, err)
		return
	}
	defer proxyAgent.Close()
	logger.Infof("Connecting to %s", realAddr)

	buf1 := make([]byte, 512)
	buf2 := make([]byte, 512)
	go func() {
		proxyAgent.Write([]byte(fmtedHeader))
		for {
			n, err := c.Read(buf1)
			proxyAgent.Write(buf1[0:n])
			if err != nil {
				break
			}
		}
	}()
	for {
		n, err := proxyAgent.Read(buf2)
		c.Write(buf2[0:n])

		if err != nil {
			break
		}
	}

	proxyAgent.Close()
	c.Close()
}

func main() {
	var err error

	// FIXME: configurable logger file
	logger = golog.New(os.Stdout)

	flag.Parse()

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
			logger.Debugf(*flDebug, "Accept return err: %s", err)
			continue
		}

		go handleTCPConnection(connTCP)
	}

}
