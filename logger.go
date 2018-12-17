package alog

import (
	"context"
	"fmt"
	"log"
	"runtime"
	"time"
)

// New returns a configured *Logger, ready to use.
//
// See the Option-returning functions in this package for configuration knobs.
func New(opts ...Option) *Logger {
	if len(opts) == 0 {
		return nil
	}
	l := Logger{}
	for _, f := range opts {
		f(&l)
	}
	if l.now == nil {
		l.now = time.Now
	}
	return &l
}

// Emitter is the interface that wraps the Emit method.
//
// Emit handles a log entry in a customized way.
type Emitter interface {
	Emit(context.Context, *Entry)
}

// EmitterFunc is an adapter to allow the use of an ordinary function as an Emitter.
type EmitterFunc func(context.Context, *Entry)

// Emit calls f(c, e).
func (f EmitterFunc) Emit(c context.Context, e *Entry) {
	f(c, e)
}

// Logger is a logging object that extracts tags from a context.Context and
// emits Entry structs.
//
// A nil *Logger is valid to call methods on.
//
// The default text format will have a newline appended if one is not present in
// the message.
type Logger struct {
	caller  bool
	emitter Emitter
	now     func() time.Time
}

// Output emits the supplied string while capturing the caller information
// "calldepth" frames back in the call stack. The value 2 is almost always what
// a caller wants. See also: runtime.Caller
func (l *Logger) Output(ctx context.Context, calldepth int, msg string) {
	// This is the bottom of the logger; everything calls this to do writes.
	// Handling nil here means everything should be able to be called on a nil
	// *Logger and not explode.
	if l == nil || l.emitter == nil {
		return
	}
	if l.now == nil {
		l.now = time.Now
	}
	e := Entry{
		Time: l.now(),
		Tags: fromContext(ctx),
		Msg:  msg,
	}

	if l.caller {
		var ok bool
		_, e.File, e.Line, ok = runtime.Caller(calldepth)
		if !ok {
			e.File = "???"
			e.Line = 0
		}
	}

	l.emitter.Emit(ctx, &e)
}

// Print calls l.Output to emit a log entry. Arguments are handled like
// fmt.Print.
func (l *Logger) Print(ctx context.Context, v ...interface{}) {
	l.Output(ctx, 2, fmt.Sprint(v...))
}

// Printf calls l.Output to emit a log entry. Arguments are handled like
// fmt.Printf.
func (l *Logger) Printf(ctx context.Context, f string, v ...interface{}) {
	l.Output(ctx, 2, fmt.Sprintf(f, v...))
}

type writerFunc func(p []byte) (int, error)

func (w writerFunc) Write(p []byte) (int, error) {
	return w(p)
}

// StdLogger returns a standard log.Logger that sends log messages to this
// logger with the provided Context. The returned log.Logger should not be
// modified.
func (l *Logger) StdLogger(ctx context.Context) *log.Logger {
	return log.New(writerFunc(func(p []byte) (int, error) {
		// standard Logger always writes a trailing \n, so remove it
		p = p[:len(p)-1]
		l.Output(ctx, 4, string(p))
		return len(p), nil
	}), "", 0)
}
