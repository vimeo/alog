package internal

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		b := &bytes.Buffer{}
		b.Grow(1024)
		return b
	},
}

// GetBuffer gets a Buffer from the pool.
func GetBuffer() *bytes.Buffer {
	return bufPool.Get().(*bytes.Buffer)
}

// PutBuffer resets a Buffer and puts it back into the pool.
func PutBuffer(b *bytes.Buffer) {
	b.Reset()
	bufPool.Put(b)
}
