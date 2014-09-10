package protocol

// Protocol:
// +.....+-----+-----+-----+-----+-----+-----+-----+-----+
// | IV  |          user name + random             | flag| encrypted by
// |     |            56bits/7bytes                | 8bit|  Master Key
// +.....+-----+-----+-----+-----+-----+-----+-----+-----+
// IV is shared for both Master Key and User Password.

// IPv4 request:                                             \
// +-----+-----+-----+-----+-----+-----+-----+               |
// |magic|      IPv4 address     |Port number|               |
// |1Byte|         4Bytes        |  2 bytes  |               |
// +-----+-----+-----+-----+-----+-----+-----+               |
//                                                           |
// IPv6 request:                                             |
// +-----+-----+-----+...........+-----+-----+-----+-----+   |
// |magic|              IPv6 address         |port number|   | encrypted by
// |1Byte|                 16Bytes           | 2 bytes   |   | User Password
// +-----+-----+-----+...........+-----+-----+-----+-----+   |
//                                                           |
// Domain request:                                           |
// +-----+-----+-----+-----+...........+-----+-----+-----+   |
// |magic| len |             Domain name                 |   |
// |1Byte|1Byte|                                         |   |
// +-----+-----+-----+-----+...........+-----+-----+-----+   /

// Request content. encrypted by User Password.

const (
	MagicByte = 0xDD
)
