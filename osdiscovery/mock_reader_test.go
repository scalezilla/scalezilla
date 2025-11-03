package osdiscovery

import (
	"errors"
)

type fakeReader struct {
	called int
}

func (f *fakeReader) Read(p []byte) (int, error) {
	f.called++
	if f.called == 1 {
		copy(p, []byte("ID=ubuntu\n"))
		return len("ID=ubuntu\n"), nil
	}
	return 0, errors.New("read error")
}
