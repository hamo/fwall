package tunnel

import (
	"io"
	"net"
	"strconv"
)

type RawSocket struct {
	c net.Conn
}

func RawSocketDial(addr string, port int) (*RawSocket, error) {
	t := addr + ":" + strconv.Itoa(int(port))

	c, err := net.Dial("tcp", t)
	if err != nil {
		return nil, err
	}

	return &RawSocket{
		c: c,
	}, nil
}

func RawSocketAccept(c net.Conn) *RawSocket {
	return &RawSocket{
		c: c,
	}
}

func (r *RawSocket) Close() {
	r.c.Close()
}

func (r *RawSocket) ReadMaster(buf []byte, full bool) (int, error) {
	if full {
		return io.ReadFull(r.c, buf)
	} else {
		return r.c.Read(buf)
	}
}

func (r *RawSocket) ReadUser(buf []byte, full bool) (int, error) {
	if full {
		return io.ReadFull(r.c, buf)
	} else {
		return r.c.Read(buf)
	}
}

func (r *RawSocket) WriteMaster(buf []byte) (int, error) {
	return r.c.Write(buf)
}

func (r *RawSocket) WriteUser(buf []byte) (int, error) {
	return r.c.Write(buf)
}
