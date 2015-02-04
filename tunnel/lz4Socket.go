package tunnel

import (
	"encoding/binary"
	"fmt"
	"io"

	"encrypt"

	"github.com/hamo/golog"
	lz4 "github.com/bkaradzic/go-lz4"
)

type LZ4SocketClient struct {
	ClientBase
}

type LZ4SocketServer struct {
	ServerBase
}

func NewLZ4SocketClient(addr string, port int, masterKey string, encryptMethod string, password string, logger *golog.GoLogger) (*LZ4SocketClient, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &LZ4SocketClient{
		ClientBase{
			addr:      addr,
			port:      port,
			crypto:    c,
			ivReady:   make(chan bool, 0),
			masterKey: c.GenKey(masterKey),
			password:  c.GenKey(password),
			logger:    logger,
		},
	}, nil
}

func NewLZ4SocketServer(masterKey string, encryptMethod string, logger *golog.GoLogger) (*LZ4SocketServer, error) {
	c, ok := encrypt.CryptoTable[encryptMethod]
	if !ok {
		return nil, fmt.Errorf("%s encrypt method is not supported.", encryptMethod)
	}

	return &LZ4SocketServer{
		ServerBase{
			crypto:    c,
			ivReady:   make(chan bool, 0),
			masterKey: c.GenKey(masterKey),
			logger:    logger,
		},
	}, nil
}


func (r *LZ4SocketClient) ReadContent(buf []byte) (int, error) {
	var n int
	var err error

	compressed_size_ib := make([]byte, 2)
	n, err = io.ReadFull(r.c, compressed_size_ib)
	if err != nil {
		return n, err
	}
	is_compressed := uint32((compressed_size_ib[0] >> 7) & 1)
	compressed_size_ib[0] &= 0x7f
	size := binary.BigEndian.Uint16(compressed_size_ib)
	stream_buf := make([]byte, size)
	n, err = io.ReadFull(r.c, stream_buf)
	if err != nil {
		return n, err
	}

	if is_compressed == 1 {
		decompressed_data, de_err := lz4.Decode(nil, stream_buf)
		if de_err != nil {
			err = de_err
		}
		n = len(decompressed_data)
		copy(buf[:n], decompressed_data[:n])
	} else {
		n = int(size)
		copy(buf[:n], stream_buf[:n])
	}

	return n, err
}

func (r *LZ4SocketServer) ReadContent(buf []byte) (int, error) {
	// call after ParseUserHeader
	var n int
	var err error

	size_ib := make([]byte, 2)
	n, err = io.ReadFull(r.c, size_ib)
	if err != nil {
		return n, err
	}
	is_compressed := uint32((size_ib[0] >> 7) & 1)
	size_ib[0] &= 0x7f
	size := binary.BigEndian.Uint16(size_ib)
	stream_buf := make([]byte, size)
	n, err = io.ReadFull(r.c, stream_buf)
	if err != nil {
		return n, err
	}

	if is_compressed == 1 {
		decompressed_data, de_err := lz4.Decode(nil, stream_buf)
		if de_err != nil {
			err = de_err
		}
		n = len(decompressed_data)
		copy(buf[:n], decompressed_data[:n])
	} else {
		n = int(size)
		copy(buf[:n], stream_buf[:n])
	}

	return n, err
}

func (r *LZ4SocketClient) WriteContent(buf []byte) (int, error) {
	compressed_buf, _ := lz4.Encode(nil, buf)
	compressed_size := uint16(len(compressed_buf))

	t_buf := make([]byte, 2)
	if compressed_size >= uint16(len(buf)) {
		binary.BigEndian.PutUint16(t_buf, uint16(len(buf)))
		t_buf = append(t_buf, buf...)
	} else {
		binary.BigEndian.PutUint16(t_buf, compressed_size)
		t_buf[0] |= 0x80
		t_buf = append(t_buf, compressed_buf...)
	}

	return r.c.Write(t_buf)
}

func (r *LZ4SocketServer) WriteContent(buf []byte) (int, error) {
	compressed_buf, _ := lz4.Encode(nil, buf)
	compressed_size := uint16(len(compressed_buf))

	t_buf := make([]byte, 2)
	if compressed_size >= uint16(len(buf)) {
		binary.BigEndian.PutUint16(t_buf, uint16(len(buf)))
		t_buf = append(t_buf, buf...)
	} else {
		binary.BigEndian.PutUint16(t_buf, compressed_size)
		t_buf[0] |= 0x80
		t_buf = append(t_buf, compressed_buf...)
	}

	return r.c.Write(t_buf)
}
