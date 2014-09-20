package main

import (
	"fmt"
	"net"

	"protocol"
	"tunnel"
)

const (
	port = ":443"
)

func handleConnection(c net.Conn) {
	r := tunnel.HttpTunnelAccept(c)

	s := protocol.NewServer(nil)

	_ = s.Accept(r)

	_, addrPort, err := s.ParseUserHeader(r)

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
