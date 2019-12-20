// Package testlog provides emitter options for tests
package testlog

import (
	"context"
	"path"
	"testing"

	"github.com/vimeo/alog/v3"
)

// DefaultLogger should be used for testing only. It returns a logger
// that will log the file path, message, and tags
func DefaultLogger(t testing.TB, opt ...Option) *alog.Logger {
	return alog.New(alog.WithCaller(), alog.WithEmitter(Emitter(t, opt...)))
}

// Emitter should be used for tests to log the file path,
// message, and tags
func Emitter(t testing.TB, opt ...Option) alog.Emitter {
	o := new(Options)
	for _, option := range opt {
		option(o)
	}

	return alog.EmitterFunc(func(ctx context.Context, e *alog.Entry) {
		t.Helper()

		if o.shortfile {
			e.File = path.Base(e.File)
		}
		t.Logf("%s %s %v", e.File, e.Msg, e.Tags)
	})
}
