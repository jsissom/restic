package backend

import (
	"io"
	"io/ioutil"
	"restic"
)

// LoadAll reads all data stored in the backend for the handle.
func LoadAll(be restic.Backend, h restic.Handle) (buf []byte, err error) {
	rd, err := be.Get(h, 0, 0)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(ioutil.Discard, rd)
		e := rd.Close()
		if err == nil {
			err = e
		}
	}()

	return ioutil.ReadAll(rd)
}

// Closer wraps an io.Reader and adds a Close() method that does nothing.
type Closer struct {
	io.Reader
}

// Close is a no-op.
func (c Closer) Close() error {
	return nil
}

// LimitedReadCloser wraps io.LimitedReader and exposes the Close() method.
type LimitedReadCloser struct {
	io.ReadCloser
	io.Reader
}

// Read reads data from the limited reader.
func (l *LimitedReadCloser) Read(p []byte) (int, error) {
	return l.Reader.Read(p)
}

// LimitReadCloser returns a new reader wraps r in an io.LimitReader, but also
// exposes the Close() method.
func LimitReadCloser(r io.ReadCloser, n int64) *LimitedReadCloser {
	return &LimitedReadCloser{ReadCloser: r, Reader: io.LimitReader(r, n)}
}
