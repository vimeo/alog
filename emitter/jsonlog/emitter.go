package jsonlog // import "github.com/vimeo/alog/emitter/jsonlog"

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/vimeo/alog/v2"
)

// DefaultLogger is a *alog.Logger with some default options
var DefaultLogger = alog.New(alog.WithEmitter(Emitter(os.Stderr, WithShortFile(), WithUTC())))

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

func jsonString(w *bytes.Buffer, s string) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(s)
	w.Truncate(w.Len() - 1)
}

// Emitter emits log messages as single lines of JSON.
//
// Logs are output to w. Every entry generates a single Write call to w, and
// calls are serialized.
func Emitter(w io.Writer, opt ...Option) alog.Emitter {
	o := new(Options)
	for _, option := range opt {
		option(o)
	}

	wOut := &out{Writer: w}

	timestampFormat := o.datefmt
	useTimestamp := o.flags&timeFlag != 0
	if !useTimestamp {
		timestampFormat = DefaultTimestampFormat
		useTimestamp = true
	}
	useTimestamp = timestampFormat != ""

	b := getBuffer()

	timestampField := o.timestampField
	if timestampField == "" {
		timestampField = DefaultTimestampField
	}
	b.Reset()
	jsonString(b, timestampField)
	timestampField = b.String()

	callerField := o.callerField
	if callerField == "" {
		callerField = DefaultCallerField
	}
	b.Reset()
	jsonString(b, callerField)
	callerField = b.String()

	messageField := o.messageField
	if messageField == "" {
		messageField = DefaultMessageField
	}
	b.Reset()
	jsonString(b, messageField)
	messageField = b.String()

	putBuffer(b)

	return alog.EmitterFunc(func(ctx context.Context, e *alog.Entry) {
		b := getBuffer()
		defer putBuffer(b)

		b.WriteByte('{')

		if useTimestamp {
			if o.flags&utcFlag != 0 {
				e.Time = e.Time.UTC()
			}
			b.WriteString(timestampField)
			b.WriteByte(':')
			jsonString(b, e.Time.Format(timestampFormat))
			b.WriteString(", ")
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
			b.WriteString(callerField)
			b.WriteByte(':')
			fb := getBuffer()
			fb.WriteString(file)
			fb.WriteByte(':')
			itoa(fb, line)
			jsonString(b, fb.String())
			putBuffer(fb)
			b.WriteString(", ")
		}

		if len(e.Tags) > 0 {
			b.WriteString(`"tags":{`)
			for i, tag := range e.Tags {
				jsonString(b, tag[0])
				b.WriteByte(':')
				jsonString(b, tag[1])
				if i < len(e.Tags)-1 {
					b.WriteString(", ")
				}
			}
			b.WriteString("}, ")
		}

		b.WriteString(messageField)
		b.WriteByte(':')
		jsonString(b, e.Msg)

		b.WriteString("}\n")

		// Writer error is swallowed, because checking errors on writing log
		// lines is my personal conception of hell.
		wOut.Write(b.Bytes())
	})
}
