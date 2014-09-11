package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"tunnel"

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
	_, err := local.ReadMaster(len)
	if err != nil {
		// FIXME
		return ""
	}

	u := make([]byte, len[0])
	_, err = local.ReadMaster(u)
	if err != nil {
		// FIXME
		return ""
	}
	s.username = string(u)

	return s.username
	// verify username and setup tunnel for user's password
}

func (s *Server) ParseUserHeader(local tunnel.Reader) (UDPconnect bool, addrPort string, err error) {
	magic := make([]byte, 1)
	_, err = local.ReadUser(magic)

	if magic[0] != MagicByte {
		return false, "", fmt.Errorf("Can not decode correctly.")
	}

	f := make([]byte, 1)
	_, err = local.ReadUser(f)

	UDP := CheckUDPFlag(f[0])

	switch {
	case CheckDomainFlag(f[0]):
		dl := make([]byte, 1)
		local.ReadUser(dl)

		d := make([]byte, dl[0])
		local.ReadUser(d)

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
	_, err = local.ReadUser(pb)
	p := binary.BigEndian.Uint16(pb)

	addrPort = addrPort + strconv.Itoa(int(p))

	return UDP, addrPort, nil
}

func (s *Server) Upstream(local tunnel.Reader, remote net.Conn) {
	// FIXME: configurable buffer size
	buf := make([]byte, 256)
	for {
		n, err := local.ReadUser(buf)
		remote.Write(buf[:n])
		if err != nil {
			break
		}
	}
}

func (s *Server) Downstream(local tunnel.Writer, remote net.Conn) {
	// FIXME: configurable buffer size
	buf := make([]byte, 256)
	for {
		n, err := remote.Read(buf)
		local.WriteUser(buf[0:n])
		if err != nil {
			break
		}
	}
}
