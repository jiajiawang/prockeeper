package prockeeper

import (
	"io"
	"io/ioutil"
	"sync"
)

// PausableWriter ...
type PausableWriter struct {
	writer  io.Writer
	discard bool
	sync.Mutex
}

// NewPausableWriter ...
func NewPausableWriter(w io.Writer) *PausableWriter {
	return &PausableWriter{
		discard: false,
		writer:  w,
	}
}

// Write ...
func (pw *PausableWriter) Write(p []byte) (n int, err error) {
	if pw.discard {
		return ioutil.Discard.Write(p)
	}

	return pw.writer.Write(p)
}

// Pause ...
func (pw *PausableWriter) Pause() {
	pw.Lock()
	defer pw.Unlock()
	pw.discard = true
}

// Resume ...
func (pw *PausableWriter) Resume() {
	pw.Lock()
	defer pw.Unlock()
	pw.discard = false
}
