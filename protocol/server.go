package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"tunnel"
	"compression"

	"github.com/hamo/golog"
)

type Server struct {
	username string

	logger *golog.GoLogger
}

func NewServer(logger *golog.GoLogger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Accept(local tunnel.Reader) string {
	len := make([]byte, 1)
	_, err := local.ReadMaster(len, true)
	if err != nil {
		// FIXME
		fmt.Printf("1")
		return ""
	}

	u := make([]byte, len[0])
	_, err = local.ReadMaster(u, true)
	if err != nil {
		// FIXME
		fmt.Printf("2")
		return ""
	}
	s.username = string(u)

	return s.username
	// verify username and setup tunnel for user's password
}

func (s *Server) ParseUserHeader(local tunnel.Reader) (UDPconnect bool, addrPort string, err error) {
	magic := make([]byte, 1)
	_, err = local.ReadUser(magic, true)

	if magic[0] != MagicByte {
		return false, "", fmt.Errorf("Can not decode correctly.")
	}

	f := make([]byte, 1)
	_, err = local.ReadUser(f, true)

	UDP := CheckUDPFlag(f[0])

	switch {
	case CheckDomainFlag(f[0]):
		dl := make([]byte, 1)
		local.ReadUser(dl, true)

		d := make([]byte, dl[0])
		local.ReadUser(d, true)

		addrPort = addrPort + string(d)

	case CheckIPv4Flag(f[0]):
		// FIXME
		fallthrough
	case CheckIPv6Flag(f[0]):
		// FIXME
		return false, "", fmt.Errorf("Not implemented")
	}

	addrPort = addrPort + ":"

	pb := make([]byte, 2)
	_, err = local.ReadUser(pb, true)
	p := binary.BigEndian.Uint16(pb)

	addrPort = addrPort + strconv.Itoa(int(p))

	return UDP, addrPort, nil
}

func (s *Server) Upstream(local tunnel.Reader, remote net.Conn) {
	buf := make([]byte, compression.BufferSize)
	for {
		n, err := local.ReadContent(buf)
		remote.Write(buf[:n])
		if err != nil {
			break
		}
	}
}

func (s *Server) Downstream(local tunnel.Writer, remote net.Conn) {
	buf := make([]byte, compression.BufferSize)
	for {
		n, err := remote.Read(buf)
		local.WriteContent(buf[0:n])
		if err != nil {
			break
		}
	}
}
