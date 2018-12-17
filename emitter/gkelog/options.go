package gkelog

import (
	"io"
)

// Options holds option values.
type Options struct {
	writer    io.Writer
	shortfile bool
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithWriter sets the writer to use for log output.
func WithWriter(w io.Writer) Option {
	return func(o *Options) { o.writer = w }
}

// WithShortFile indicates to only use the filename instead of the full file
// path for sourceLocation.
func WithShortFile() Option {
	return func(o *Options) { o.shortfile = true }
}
