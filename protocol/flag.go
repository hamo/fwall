package protocol

// Bit 0: If set, indicates that this is an UDP request, or, it is a TCP request.

// Bit 2  Bit 1
//   0      0   IPv4 request
//   0      1   IPv6 request
//   1      0   Domain name request
//   1      1   reserved

// Bit 3: If set, the content will be compressed.
// Bit 4: If set, the content will be encrypted.
// Warning: if only one bit is enough for encryption? it need be extended sometime.

// Bit 7 - Bit 5: reserved

const (
	connectTypeFlag byte = 1 << iota
	inetFamilyFlag  byte = 1 << iota
	addressTypeFlag byte = 1 << iota
	compressFlag    byte = 1 << iota
	encryptionFlag  byte = 1 << iota
)

func NewFlag() *byte {
	var f byte
	return &f
}

func SetUDPFlag(f *byte) {
	*f |= connectTypeFlag
}

func SetTCPFlag(f *byte) {
	*f &= ^connectTypeFlag
}

func CheckUDPFlag(f byte) bool {
	return (f & connectTypeFlag) == 1
}

func CheckTCPFlag(f byte) bool {
	return (f & connectTypeFlag) == 0
}

func SetDomainFlag(f *byte) {
	*f &= ^inetFamilyFlag
	*f |= addressTypeFlag
}

func SetIPv4Flag(f *byte) {
	*f &= ^addressTypeFlag
	*f &= ^inetFamilyFlag
}

func SetIPv6Flag(f *byte) {
	*f &= ^addressTypeFlag
	*f |= inetFamilyFlag
}

func CheckDomainFlag(f byte) bool {
	return f&addressTypeFlag == addressTypeFlag &&
		f&inetFamilyFlag == 0
}

func CheckIPv4Flag(f byte) bool {
	return f&addressTypeFlag == 0 &&
		f&inetFamilyFlag == 0
}

func CheckIPv6Flag(f byte) bool {
	return f&addressTypeFlag == 0 &&
		f&inetFamilyFlag == inetFamilyFlag
}

func SetCompressFlag(f *byte) {
	*f |= compressFlag
}

func UnsetCompressFlag(f *byte) {
	*f &= ^compressFlag
}

func SetEncryptionFlag(f *byte) {
	*f |= encryptionFlag
}

func UnsetEncryptionFlag(f *byte) {
	*f &= ^encryptionFlag
}

func CheckCompressFlag(f byte) bool {
	return f&compressFlag == 1
}

func CheckEncryptionFlag(f byte) bool {
	return f&encryptionFlag == 1
}
