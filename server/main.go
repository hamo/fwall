package main

import (
	"net"
)

const (
	port = ":3389"
)

func handleConnection(c net.Conn) {

}

func main() {

	lnTCP, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}

	for {
		conn, err := lnTCP.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}
