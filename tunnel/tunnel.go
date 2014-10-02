package tunnel

import ()

type Reader interface {
	// Read IV and Master header
	ReadMaster(p []byte, full bool) (int, error)

	// Read User header
	ReadUser(p []byte, full bool) (int, error)

	// Read Content data
	ReadContent(p []byte) (int, error)
}

type Writer interface {
	// Write IV and Master header
	WriteMaster(p []byte) (n int, err error)

	// Write User header
	WriteUser(p []byte) (n int, err error)

	// Write Content data
	WriteContent(p []byte) (n int, err error)
}
