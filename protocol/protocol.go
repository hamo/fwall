package protocol

// Protocol: 
// +-----+-----+-----+-----+-----+-----+-----+-----+
// |           user name + random            | flag| encrypted by 
// |              56bits/7bytes              | 8bit|  Master Key
// +-----+-----+-----+-----+-----+-----+-----+-----+

// IPv4 request:                                       \
// +-----+-----+-----+-----+-----+-----+               |
// |     IPv4 address      |port number|               |
// |        4Bytes         |  2 bytes  |               |
// +-----+-----+-----+-----+-----+-----+               |
                                                       |
// IPv6 request:                                       |
// +-----+-----+...........+-----+-----+-----+-----+   |
// |           IPv6 address            |port number|   | encrypted by
// |              16Bytes              |  2 bytes  |   | User Password
// +-----+-----+...........+-----+-----+-----+-----+   |
                                                       |
// Domain request:                                     |
// +-----+-----+-----+...........+-----+-----+-----+   |
// | len |            Domain name                  |   |
// |1Byte|                                         |   |
// +-----+-----+-----+...........+-----+-----+-----+   /

// Request content. encrypted by User Password.
