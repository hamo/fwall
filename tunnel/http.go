package tunnel

import (
	"net"
	"strconv"
	"io"
)

type HttpTunnel struct {
	c net.Conn
}

func HttpTunnelDial(addr string, port int) (*HttpTunnel, error) {
	t := addr + ":" + strconv.Itoa(int(port))

	c, err := net.Dial("tcp", t)
	if err != nil {
		return nil, err
	}

	return &HttpTunnel{
		c: c,
	}, nil
}

func HttpTunnelAccept(c net.Conn) *HttpTunnel {
	return &HttpTunnel{
		c: c,
	}
}

func (r *HttpTunnel) Close() {
	r.c.Close()
}

func (r *HttpTunnel) ReadMaster(buf []byte, full bool) (int, error) {
	if full {
		return io.ReadFull(r.c, buf)
	} else {
		return r.c.Read(buf)
	}
}

func (r *HttpTunnel) ReadUser(buf []byte, full bool) (int, error) {
	if full {
		return io.ReadFull(r.c, buf)
	} else {
		return r.c.Read(buf)
	}
}
func (r *HttpTunnel) WriteMaster(buf []byte) (int, error) {
	return r.c.Write(buf)
}

func (r *HttpTunnel) WriteUser(buf []byte) (int, error) {
	return r.c.Write(buf)
}
