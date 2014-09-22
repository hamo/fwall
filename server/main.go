package main

import (
	"fmt"
	"net"
	"os"

	"protocol"
	"tunnel"

	"github.com/hamo/golog"
)

const (
	port = ":443"
)

func handleConnection(c net.Conn) {
	logger := golog.New(os.Stdout)
	logger.SetDebug(true)

	r, err := tunnel.NewRawSocket("", 443, "server", "foobar", "aes-256-cfb", "barfoo", logger)

	r.Accept(c)

	s := protocol.NewServer(nil)

	_ = s.Accept(r)

	_, addrPort, err := s.ParseUserHeader(r)
	if err != nil {
		logger.Errorf("ParseUserHeader failed: %v", err)
		return
	}

	realServer, err := net.Dial("tcp", addrPort)

	if err != nil {
		fmt.Printf("err: %s", err)
		return
	}

	go s.Upstream(r, realServer)
	s.Downstream(r, realServer)

	realServer.Close()
	c.Close()
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
