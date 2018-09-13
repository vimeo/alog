package alog

import (
	"io"
	"sync"
)

const (
	fileFlag = 1 << iota
	shortfileFlag
	timeFlag
	utcFlag
)

// Option sets an option on a Logger.
//
// Options are applied in the order specified, meaning if both To and
// WithEmitter are supplied in a call to New, the last one wins.
type Option func(*Logger)

// To configures a Logger to output log lines in a human-readable format to the
// supplied io.Writer.
//
// It's equivalent to
//
//	var w io.Writer
//	l := New(func(l *Logger) {
//		WithEmitter(l.EmitText(w))(l)
//	})
//
// with additional guarentees that every entry generates a single Write call,
// and calls are serialized.
func To(w io.Writer) Option {
	return func(l *Logger) { l.emitter = l.EmitText(&out{Writer: w}) }
}

// out is a wrapper to guarentee serialized access to the inner Writer.
type out struct {
	sync.Mutex
	io.Writer
}

func (o *out) Write(b []byte) (int, error) {
	o.Lock()
	n, err := o.Writer.Write(b)
	o.Unlock()
	return n, err
}

// WithDateFormat sets the string format for timestamps using a layout string
// like the time package would take.
//
// This only effects the default text format.
func WithDateFormat(layout string) Option {
	return func(l *Logger) {
		l.datefmt = layout
		l.flags |= timeFlag
	}
}

// WithPrefix adds a set prefix to all lines.
//
// This only effects the default text format.
func WithPrefix(prefix string) Option {
	return func(l *Logger) { l.prefix = prefix }
}

// WithFile collects call information on each log line, like the log
// package's Llongfile flag.
func WithFile() Option {
	return func(l *Logger) { l.flags |= fileFlag }
}

// WithShortFile is like WithFile, but only prints the file name
// instead of the entire path.
//
// This only effects the default text format.
func WithShortFile() Option {
	return func(l *Logger) { l.flags |= fileFlag | shortfileFlag }
}

// WithUTC sets timestamps to UTC.
//
// This only effects the default text format.
func WithUTC() Option {
	return func(l *Logger) { l.flags |= utcFlag }
}

// WithEmitter configures the logger to call f every time it needs to emit a log
// line.
//
// Calls to f are not synchronized.
func WithEmitter(e Emitter) Option {
	return func(l *Logger) { l.emitter = e }
}
