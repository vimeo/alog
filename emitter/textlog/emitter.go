package textlog // import "github.com/vimeo/alog/emitter/textlog"

import (
	"bytes"
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/vimeo/alog/v2"
)

// Default is an alog.Emitter with some default options
var Default = alog.New(alog.WithEmitter(Emitter(os.Stderr, WithShortFile(), WithDateFormat(time.RFC3339), WithUTC())))

// out is a wrapper to guarantee serialized access to the inner Writer.
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

// Emitter emits log messages as plain text.
//
// Logs are output to w. The format is determined by l.
// Every entry generates a single Write call to w, and calls are serialized.
func Emitter(w io.Writer, opt ...Option) alog.Emitter {
	o := new(Options)
	for _, option := range opt {
		option(o)
	}
	wOut := &out{Writer: w}
	return alog.EmitterFunc(func(ctx context.Context, e *alog.Entry) {
		m := getBuffer()
		defer putBuffer(m)
		m.WriteString(o.prefix)

		// Quick-n-dirty custom formatting code
		if o.flags&timeFlag != 0 {
			if o.flags&utcFlag != 0 {
				e.Time = e.Time.UTC()
			}
			m.WriteString(e.Time.Format(o.datefmt))
			m.WriteByte(' ')
		}
		if o.flags&fileFlag != 0 && e.File != "" {
			file := e.File
			line := e.Line
			if o.flags&shortfileFlag != 0 {
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
			m.WriteString(": ")
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
		wOut.Write(m.Bytes())
	})
}
