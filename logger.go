package alog

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
	return &l
}

// Logger is a logging object that extracts tags from a context.Context and
// emits Entry structs.
//
// A nil *Logger is valid to call methods on.
//
// The default text format will have a newline appended if one is not present in
// the message.
type Logger struct {
	prefix  string
	datefmt string
	flags   uint

	emit func(*Entry)
}

// Output emits the supplied string while capturing the caller information
// "calldepth" frames back in the call stack. The value 2 is almost always what
// a caller wants. See also: runtime.Caller
func (l *Logger) Output(ctx context.Context, calldepth int, msg string) {
	// This is the bottom of the logger; everything calls this to do writes.
	// Handling nil here means everything should be able to be called on a nil
	// *Logger and not explode.
	if l == nil || l.emit == nil {
		return
	}
	e := Entry{
		Time: time.Now(),
		Tags: fromContext(ctx),
		Msg:  msg,
	}
	if l.flags&fileFlag != 0 {
		var ok bool
		_, e.File, e.Line, ok = runtime.Caller(calldepth)
		if !ok {
			e.File = "???"
			e.Line = 0
		}
	}

	l.emit(&e)
}

// EmitText is the default text format emitter.
//
// It closes over the supplied io.Writer and returns a function suitable to
// pass to WithEmitter.
func (l *Logger) EmitText(w io.Writer) func(e *Entry) {
	return func(e *Entry) {
		m := getBuffer()
		defer putBuffer(m)
		m.WriteString(l.prefix)

		// Quick-n-dirty custom formatting code
		if l.flags&fileFlag != 0 {
			file := e.File
			line := e.Line
			if l.flags&shortfileFlag != 0 {
				for i := len(e.File) - 1; i > 0; i-- {
					if file[i] == '/' {
						file = file[i+1:]
						break
					}
				}
			}
			m.WriteString(file)
			m.WriteByte(':')
			itoa(m, line)
			m.WriteByte(' ')
		}

		if l.flags&timeFlag != 0 {
			if l.flags&utcFlag != 0 {
				e.Time = e.Time.UTC()
			}
			m.WriteString(e.Time.Format(l.datefmt))
			m.WriteByte(' ')
		}
		if t := e.Tags; len(t) != 0 {
			m.WriteByte('[')
			for i, p := range t {
				if i != 0 {
					m.WriteByte(' ')
				}
				m.WriteString(p[0])
				m.WriteByte('=')
				m.WriteString(p[1])
			}
			m.WriteString("] ")
		}
		m.WriteString(e.Msg)
		if m.Bytes()[m.Len()-1] != '\n' {
			m.WriteByte('\n')
		}
		// Writer error is swallowed, because checking errors on writing log
		// lines is my personal conception of hell.
		w.Write(m.Bytes())
	}
}

func itoa(w *bytes.Buffer, i int) {
	buf := make([]byte, 16)
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
