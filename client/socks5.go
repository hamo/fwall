package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

var (
	handShakeProtoVersion = 0
	handShakeAnswer       = []byte{0x05, 0x00}
	reqAnswer             = []byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43}
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
func parseReq(c net.Conn) (byte, byte, []byte, uint16, error) {
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

	var port uint16
	binary.Read(r, binary.BigEndian, &port)

	// commandCode, addressType, address, port, err
	return commandCode, addressType, address.Bytes(), port, nil
}
