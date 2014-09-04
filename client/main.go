package main

import (
	"net"
	"os"
	"strconv"

	"github.com/hamo/golog"
)

const (
	version = "0.1"
)

const (
	port           = ":1081"
	remote_port    = ":443"
	remote_address = "128.199.153.182"
)

var (
	debug  = true
	logger *golog.GoLogger
)

func fmtHeader(remoteAddr []byte) []byte {
	lengthByte := []byte{byte(len(remoteAddr))}
	fmtedHeader := append(lengthByte, remoteAddr...)
	return fmtedHeader
}

func handleTCPConnection(c net.Conn) {
	err := handShake(c)

	if err != nil {
		logger.Debugf(debug, "handShake err: %s", err)
		c.Close()
		return
	}

	commandCode, addressType, address, port, err := parseReq(c)
	if err != nil {
		logger.Fatalf("parseReq failed: %s", err)
	}

	logger.Debugf(debug, "commandCode: %d\n", commandCode)
	logger.Debugf(debug, "addressType: %d\n", addressType)
	logger.Debugf(debug, "port: %d\n", port)

	c.Write(reqAnswer)

	// Move to server
	//	realAddr := string(address)
	var realAddr string
	if addressType == 0x03 {
		realAddr = string((address.Bytes()[1:]))
	}

	realAddr = realAddr + ":" + strconv.Itoa(int(port))
	remoteAddr := remote_address + remote_port
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
	// FIXME: configurable logger file
	logger = golog.New(os.Stdout)

	logger.Infof("fwall started. Version: %s", version)

	lnTCP, err := net.Listen("tcp", port)
	if err != nil {
		logger.Fatalf("Listen to socks5 port failed: %s", err)
	}
	defer lnTCP.Close()

	for {
		connTCP, err := lnTCP.Accept()
		if err != nil {
			logger.Debugf(debug, "Accept return err: %s", err)
			continue
		}

		go handleTCPConnection(connTCP)
	}

}
