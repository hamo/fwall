package protocol

// Protocol:
// +.....+-----+-----+-----+.....+-----+-----+
// | IV  | len |          username           | encrypted by
// |     |1Byte|                             |  Master Key
// +.....+-----+-----+-----+.....+-----+-----+
// IV is shared for both Master Key and User Password.

// IPv4 request:                                                   \
// +-----+-----+-----+-----+-----+-----+-----+-----+               |
// |magic| flag|      IPv4 address     |Port number|               |
// |1Byte| 8bit|         4Bytes        |  2 bytes  |               |
// +-----+-----+-----+-----+-----+-----+-----+-----+               |
//                                                                 |
// IPv6 request:                                                   |
// +-----+-----+-----+-----+...........+-----+-----+-----+-----+   |
// |magic| flag|           IPv6 address            |port number|   | encrypted by
// |1Byte| 8bit|              16Bytes              | 2 bytes   |   | User Password
// +-----+-----+-----+-----+...........+-----+-----+-----+-----+   |
//                                                                 |
// Domain request:                                                 |
// +-----+-----+-----+-----+-----+...........+-----+-----+-----+   |
// |magic| flag| len |            Domain name                  |   |
// |1Byte| 8bit|1Byte|                                         |   |
// +-----+-----+-----+-----+-----+...........+-----+-----+-----+   /

// Request content. encrypted by User Password.

const (
	MagicByte = 0xDD
)
