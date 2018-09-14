package alog

import (
	"time"
)

// Option sets an option on a Logger.
//
// Options are applied in the order specified, meaning if both To and
// WithEmitter are supplied in a call to New, the last one wins.
type Option func(*Logger)

// WithEmitter configures the logger to call e.Emit() every time it needs to
// emit a log line.
//
// Calls are not synchronized.
func WithEmitter(e Emitter) Option {
	return func(l *Logger) { l.emitter = e }
}

// WithCaller configures the logger to include the caller information in each
// log entry.
func WithCaller() Option {
	return func(l *Logger) { l.caller = true }
}

// OverrideTimestamp sets the function that will be used to get the current
// time for each log entry.
//
// This is primarily meant to be used for testing custom emitters.
// For example: OverrideTimestamp(func() time.Time { return time.Time{} })
func OverrideTimestamp(f func() time.Time) Option {
	return func(l *Logger) { l.now = f }
}
