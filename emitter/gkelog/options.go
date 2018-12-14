package gkelog

import (
	"io"
)

// Options holds option values.
type Options struct {
	writer  io.Writer
	project string
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithWriter sets the writer to use for log output.
func WithWriter(w io.Writer) Option {
	return func(o *Options) { o.writer = w }
}
