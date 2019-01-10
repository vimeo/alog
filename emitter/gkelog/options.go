package gkelog

import (
	"io"
)

// Options holds option values.
type Options struct {
	reqWriter io.Writer
	appWriter io.Writer
	shortfile bool
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithWriter sets the writer to use for log output.
func WithWriter(w io.Writer) Option {
	return func(o *Options) {
		o.reqWriter = w
		o.appWriter = w
	}
}

// WithWriters sets the writers to use for log output.
//
// req is used for log entries that have a Request.
// app is used for all other log entries.
// In GKE the typical usage would be WithWriters(os.Stdout, os.Stderr).
func WithWriters(req io.Writer, app io.Writer) Option {
	return func(o *Options) {
		o.reqWriter = req
		o.appWriter = app
	}
}

// WithShortFile indicates to only use the filename instead of the full file
// path for sourceLocation.
func WithShortFile() Option {
	return func(o *Options) { o.shortfile = true }
}
