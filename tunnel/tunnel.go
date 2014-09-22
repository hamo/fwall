package tunnel

import ()

type Reader interface {
	// Read IV and Master header
	ReadMaster(p []byte, full bool) (int, error)

	// Read User data
	ReadUser(p []byte, full bool) (int, error)
}

type Writer interface {
	// Write IV and Master header
	WriteMaster(p []byte) (n int, err error)

	// Write User data
	WriteUser(p []byte) (n int, err error)
}
