package internal

import (
	"bytes"
	"io"
	"sync"
)

// SerializedWriter is a wrapper to guarantee serialized access to the inner Writer.
type SerializedWriter struct {
	sync.Mutex
	io.Writer
}

// NewSerializedWriter creates a new SerializedWriter using w as the inner Writer.
func NewSerializedWriter(w io.Writer) *SerializedWriter {
	return &SerializedWriter{Writer: w}
}

// Write writes to the inner Writer. Concurrent calls to Write will block so
// that only one at a time is writing to the inner Writer.
func (o *SerializedWriter) Write(b []byte) (int, error) {
	o.Lock()
	n, err := o.Writer.Write(b)
	o.Unlock()
	return n, err
}

// Itoa writes the string representation of an int to a Buffer.
func Itoa(w *bytes.Buffer, i uint) {
	buf := make([]byte, 20)
	p := len(buf) - 1
	for i >= 10 {
		q := i / 10
		buf[p] = byte('0' + i - q*10)
		p--
		i = q
	}
	buf[p] = byte('0' + i)
	w.Write(buf[p:])
}
