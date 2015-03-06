package tunnel

import (
	"crypto/aes"
	"github.com/hamo/fwall/encrypt"
	"io"
	"net"
	"strconv"
)

var eTunnel encrypt.Encryption

type HttpTunnel struct {
	c net.Conn
}

func HttpTunnelDial(addr string, port int) (*HttpTunnel, error) {
	t := addr + ":" + strconv.Itoa(int(port))

	c, err := net.Dial("tcp", t)
	if err != nil {
		return nil, err
	}

	// FIXME: Seems it is the best place to init the cipher. Fix it
	// if you have a better choice. ;)

	// Let's using fixed iv and key first.
	iv := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")[:aes.BlockSize]
	Key := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")[0:16]
	eTunnel = encrypt.NewAesStream(iv, Key)

	return &HttpTunnel{
		c: c,
	}, nil
}

func HttpTunnelAccept(c net.Conn) *HttpTunnel {
	iv := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")[:aes.BlockSize]
	Key := []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")[0:16]
	eTunnel = encrypt.NewAesStream(iv, Key)
	return &HttpTunnel{
		c: c,
	}
}

func (r *HttpTunnel) Close() {
	r.c.Close()
}

func (r *HttpTunnel) ReadMaster(buf []byte, full bool) (int, error) {
	encryptedBuf := make([]byte, len(buf))

	var len int
	var err error
	if full {
		len, err = io.ReadFull(r.c, encryptedBuf)
	} else {
		len, err = r.c.Read(encryptedBuf)
	}
	uncryptedBuf := eTunnel.Decrypt(encryptedBuf)
	copy(buf, uncryptedBuf)

	return len, err
}

func (r *HttpTunnel) ReadUser(buf []byte, full bool) (int, error) {
	if full {
		return io.ReadFull(r.c, buf)
	} else {
		return r.c.Read(buf)
	}
}
func (r *HttpTunnel) WriteMaster(buf []byte) (int, error) {

	encryptedBuf := eTunnel.Encrypt(buf)

	return r.c.Write(encryptedBuf)
}

func (r *HttpTunnel) WriteUser(buf []byte) (int, error) {
	return r.c.Write(buf)
}
