package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	port = ":1081"
	remote_port = ":443"
	remote_address = "128.199.153.182"
)

const (
	handShakeProtoVersion = 0
)

var (
	handShakeAnswer = []byte{0x05, 0x00}
	reqAnswer       = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43}
)

func handShake(c net.Conn) error {
	hsr := make([]byte, 3)
	_, err := io.ReadFull(c, hsr)

	if err != nil {
		return fmt.Errorf("HandShake error")
	}

	protoVersion := hsr[handShakeProtoVersion]

	if protoVersion != 0x05 {
		return fmt.Errorf("Protocol mismatch: %d", protoVersion)
	}

	c.Write(handShakeAnswer)
	return nil
}

// commandCode, addressType, address, port, err
func parseReq(c net.Conn) (byte, byte, *bytes.Buffer, int16, error) {
	r := bufio.NewReader(c)

	protoVersion, err := r.ReadByte()
	if err != nil {
		return 0, 0, nil, 0, err
	}

	if protoVersion != 0x05 {
		return 0, 0, nil, 0, fmt.Errorf("Protocol mismatch: %d", protoVersion)
	}

	commandCode, err := r.ReadByte()
	if err != nil {
		return 0, 0, nil, 0, err
	}

	_, err = r.ReadByte()
	if err != nil {
		return 0, 0, nil, 0, err
	}

	addressType, err := r.ReadByte()
	if err != nil {
		return 0, 0, nil, 0, err
	}

	address := new(bytes.Buffer)

	switch addressType {
	case 0x01:
		for i := 0; i < 4; i++ {
			c, _ := r.ReadByte()
			address.WriteByte(c)
		}
	case 0x03:
		length, _ := r.ReadByte()
		address.WriteByte(length)
		for i := 0; i < int(length); i++ {
			c, _ := r.ReadByte()
			address.WriteByte(c)
		}
	case 0x04:
		for i := 0; i < 16; i++ {
			c, _ := r.ReadByte()
			address.WriteByte(c)
		}
	default:
		return 0, 0, nil, 0, fmt.Errorf("Address Type error")
	}

	var port int16
	binary.Read(r, binary.BigEndian, &port)

	// commandCode, addressType, address, port, err
	return commandCode, addressType, address, port, nil
}
func fmtHeader(remoteAddr []byte) ([]byte){
	lengthByte := []byte{byte(len(remoteAddr))}
	fmtedHeader := append(lengthByte, remoteAddr...)
	return fmtedHeader
}

func handleTCPConnection(c net.Conn) {
	err := handShake(c)

	if err != nil {
		fmt.Printf("DEBUG: handShake err: %s", err)
		c.Close()
		return
	}

	commandCode, addressType, address, port, err := parseReq(c)
	if err != nil {
		panic(err)
	}

	fmt.Printf("commandCode: %d\n", commandCode)
	fmt.Printf("addressType: %d\n", addressType)
	// fmt.Printf("address: %s\n", string(address.Bytes()))
	fmt.Printf("port: %d\n", port)

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
	fmt.Printf("realAddr: %s, \t remoteAddr: %s \t formatd addr: %s\n", realAddr, remoteAddr, string(fmtedHeader))

	proxyAgent, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		panic(err)
	}
	fmt.Printf("succeed dial to %s\n", remoteAddr)

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

	proxyAgent.Read(buf2)
	c.Write(buf2)
}

func main() {
	lnTCP, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	defer lnTCP.Close()

	for {
		connTCP, err := lnTCP.Accept()
		if err != nil {
			continue
		}

		go handleTCPConnection(connTCP)
	}

}
