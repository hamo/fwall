package main

import (
	"fmt"
	"net"
)

const (
	port = ":443"
)

func handleConnection(c net.Conn) {
	length := make([]byte, 1)
	c.Read(length)
	lnum := int(length[0])
	realAddr := make([]byte, lnum)
	c.Read(realAddr)
	fmt.Printf("RealAddr: %s\n", string(realAddr))
	realServer, err := net.Dial("tcp", string(realAddr))
	if err != nil {
		return
	}
	defer realServer.Close()

	buf1 := make([]byte, 512)
	buf2 := make([]byte, 512)

	go func() {
		for {
			n, err := c.Read(buf1)
			realServer.Write(buf1[0:n])
			if err != nil {
				break
			}
		}
	}()
	for {
		n, err := realServer.Read(buf2)
		c.Write(buf2[0:n])

		if err != nil {
			break
		}
	}

	realServer.Close()
	c.Close()
}
func repeater(c net.Conn) {

}
func main() {

	lnTCP, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	defer lnTCP.Close()

	for {
		conn, err := lnTCP.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}
