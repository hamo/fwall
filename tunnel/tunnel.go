package tunnel

import (
	"net"
	"time"
)

// Protocol: 
// +----------------+----------------+----------------+
// |                |    Byte Count                   |
// |         8 bits |                      16 bits    |
// +----------------+----------------+----------------+

type tunnel interface {
	pipe(from, to net.Conn, timeout time.Duration)
}
