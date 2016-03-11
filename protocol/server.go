package protocol

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"

	"github.com/hamo/fwall/compression"
	"github.com/hamo/fwall/tunnel"

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
		return ""
	}

	u := make([]byte, len[0])
	_, err = local.ReadMaster(u, true)
	if err != nil {
		return ""
	}
	s.username = string(u)

	return s.username
	// verify username and setup tunnel for user's password
}

func (s *Server) ParseUserHeader(local tunnel.Reader) (UDPconnect bool, addrPort string, err error) {
	magic := make([]byte, 1)
	_, err = local.ReadUser(magic, true)
	if err != nil {
		return false, "", fmt.Errorf("Failed to read magic byte.")
	}
	if magic[0] != MagicByte {
		return false, "", fmt.Errorf("Can not decode correctly.")
	}

	f := make([]byte, 1)
	_, err = local.ReadUser(f, true)
	if err != nil {
		return false, "", fmt.Errorf("Failed to read flag.")
	}
	UDP := CheckUDPFlag(f[0])

	switch {
	case CheckDomainFlag(f[0]):
		dl := make([]byte, 2)
		local.ReadUser(dl, true)
		dLen := binary.BigEndian.Uint16(dl)

		d := make([]byte, dLen)
		local.ReadUser(d, true)

		addrPort = addrPort + string(d)

	case CheckIPv4Flag(f[0]):
		v4 := make([]byte, 4)
		local.ReadUser(v4, true)
		addrPort = addrPort + net.IP(v4).String()

	case CheckIPv6Flag(f[0]):
		v6 := make([]byte, 16)
		local.ReadUser(v6, true)
		addrPort = "[" + net.IP(v6).String() + "]"
	}

	addrPort = addrPort + ":"

	pb := make([]byte, 2)
	_, err = local.ReadUser(pb, true)
	if err != nil {
		return false, "", fmt.Errorf("Failed to read port.")
	}
	p := binary.BigEndian.Uint16(pb)

	addrPort = addrPort + strconv.Itoa(int(p))

	return UDP, addrPort, nil
}

func (s *Server) Upstream(local tunnel.Reader, remote net.Conn) {
	buf := make([]byte, compression.BufferSize)
	for {
		n, err := local.ReadContent(buf)
		m, err2 := remote.Write(buf[:n])
		if err != nil || err2 != nil || m != n {
			break
		}
	}
}

func (s *Server) Downstream(local tunnel.Writer, remote net.Conn) {
	buf := make([]byte, compression.BufferSize)
	for {
		n, err := remote.Read(buf)
		m, err2 := local.WriteContent(buf[0:n])
		if err != nil || err2 != nil || m != n {
			break
		}
	}
}
