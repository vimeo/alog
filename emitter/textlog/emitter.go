package textlog

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/vimeo/alog/v3"
	"github.com/vimeo/alog/v3/emitter/internal"
)

// Default is an alog.Emitter with some default options
var Default = alog.New(alog.WithEmitter(Emitter(os.Stderr, WithShortFile(), WithDateFormat(time.RFC3339), WithUTC())))

// Emitter emits log messages as plain text.
//
// Logs are output to w. The format is determined by l.
// Every entry generates a single Write call to w, and calls are serialized.
func Emitter(w io.Writer, opt ...Option) alog.Emitter {
	o := new(Options)
	for _, option := range opt {
		option(o)
	}
	wOut := internal.NewSerializedWriter(w)
	return alog.EmitterFunc(func(ctx context.Context, e *alog.Entry) {
		m := internal.GetBuffer()
		defer internal.PutBuffer(m)
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
			line := uint(e.Line)
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
			internal.Itoa(m, line)
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
