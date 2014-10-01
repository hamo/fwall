package protocol

import (
	"encoding/binary"

	"net"
	"tunnel"
	"compression"

	"github.com/hamo/golog"
)

type Client struct {
	username string

	addressType byte
	address     []byte
	port        uint16

	logger *golog.GoLogger
}

func NewClient(username string, addressType byte, address []byte, port uint16, logger *golog.GoLogger) *Client {
	return &Client{
		username:    username,
		addressType: addressType,
		address:     address,
		port:        port,
		logger:      logger,
	}
}

func (c *Client) Upstream(client net.Conn, server tunnel.Writer) {
	l := len(c.username)
	bl := byte(l & 0xFF)
	nu := c.username[:bl]
	if c.username != nu {
		c.logger.Warningf("%s is too long. Trunc to %s", c.username, nu)
	}
	c.username = nu

	mh := make([]byte, 1+bl)
	mh[0] = bl
	copy(mh[1:1+bl], c.username)

	server.WriteMaster(mh)

	uh := make([]byte, 1)
	uh[0] = MagicByte

	// FIXME: handle TCP/UDP flag
	f := NewFlag()

	switch c.addressType {
	case 0x01: //IPv4
		SetIPv4Flag(f)
	case 0x03: //Domain
		SetDomainFlag(f)
	case 0x04: //IPv6
		SetIPv6Flag(f)
	}

	uh = append(uh, *f)

	switch c.addressType {
	case 0x01: //IPv4
		fallthrough
	case 0x04: //IPv6
		uh = append(uh, c.address...)
	case 0x03: //Domain
		ali := len(c.address)
		alb := byte(ali & 0xFF)
		uh = append(uh, alb)
		uh = append(uh, c.address...)
	}

	p := make([]byte, 2)
	binary.BigEndian.PutUint16(p, c.port)

	uh = append(uh, p...)

	server.WriteUser(uh)

	buf := make([]byte, compression.BufferSize)

	for {
		n, err := client.Read(buf)
		server.WriteContent(buf[:n])
		if err != nil {
			break
		}
	}
}

func (c *Client) Downstream(client net.Conn, server tunnel.Reader) {
	buf := make([]byte, compression.BufferSize)
	for {
		n, err := server.ReadContent(buf)
		client.Write(buf[:n])
		if err != nil {
			break
		}
	}
}
