package jsonlog

import (
	"io"
)

const (
	fileFlag = 1 << iota
	shortfileFlag
	timeFlag
	utcFlag
)

const (
	// DefaultTimestampFormat is the default value used for timestamps.
	// This is the same as time.RFC3339Nano except that it has zero-padding of
	// the fractional seconds.
	DefaultTimestampFormat = "2006-01-02T15:04:05.000000000Z07:00"

	// DefaultTimestampField is the default field name used for the timestamp
	// value. It is used if WithTimestampField is not specified.
	DefaultTimestampField = "timestamp"

	// DefaultCallerField is the default field name used for the caller
	// information. It is used if WithCallerField is not specified.
	DefaultCallerField = "caller"

	// DefaultMessageField is the default field name used for the log message
	// It is used if WithMessageField is not specified.
	DefaultMessageField = "message"
)

// Options holds option values.
type Options struct {
	timestampField string
	callerField    string
	messageField   string
	datefmt        string
	flags          uint
	writer         io.Writer
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithDateFormat sets the string format for timestamps using a layout string
// like the time package would take.
func WithDateFormat(layout string) Option {
	return func(o *Options) {
		o.datefmt = layout
		o.flags |= timeFlag
	}
}

// WithUTC sets timestamps to UTC.
func WithUTC() Option {
	return func(o *Options) { o.flags |= utcFlag }
}

// WithTimestampField overrides the JSON field used for the timestamp.
//
// If this option is not specified, DefaultTimestampField will be used.
func WithTimestampField(field string) Option {
	return func(o *Options) { o.timestampField = field }
}

// WithFile collects call information on each log line, like the log
// package's Llongfile flag.
//
// The alog.WithCaller() option also needs to be used when creating the Logger
// in order to have the file and line information added to the log entries.
func WithFile() Option {
	return func(o *Options) { o.flags |= fileFlag }
}

// WithShortFile is like WithFile, but only prints the file name
// instead of the entire path.
//
// The alog.WithCaller() option also needs to be used when creating the Logger
// in order to have the file and line information added to the log entries.
func WithShortFile() Option {
	return func(o *Options) { o.flags |= fileFlag | shortfileFlag }
}

// WithCallerField overrides the JSON field used for the caller file and line
// number information.
//
// If this option is not specified, DefaultCallerField will be used.
func WithCallerField(field string) Option {
	return func(o *Options) { o.callerField = field }
}

// WithMessageField overrides the JSON field used for the log message.
//
// If this option is not specified, DefaultMessageField will be used.
func WithMessageField(field string) Option {
	return func(o *Options) { o.messageField = field }
}

// WithWriter sets the writer to use for log output.
func WithWriter(w io.Writer) Option {
	return func(o *Options) { o.writer = w }
}
