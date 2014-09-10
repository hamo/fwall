package tunnel

import (
	"net"
	"time"
)

type tunnel interface {
	pipe(from, to net.Conn, timeout time.Duration)
}
