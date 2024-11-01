package leveled

import (
	"bytes"
	"context"
	"testing"

	"github.com/vimeo/alog/v3"
	"github.com/vimeo/alog/v3/emitter/textlog"
)

func TestLogger(t *testing.T) {
	b := &bytes.Buffer{}
	l := Default(alog.New(alog.WithEmitter(textlog.Emitter(b))))

	ctx := context.Background()
	ctx = alog.AddTags(ctx, "key", "value")
	l.Info(ctx, "")
	const want = `[key=value level=info] ` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got: %#q, want: %#q", got, want)
	}
}

func TestLogLevels(t *testing.T) {
	b := &bytes.Buffer{}
	l := Filtered(alog.New(alog.WithEmitter(textlog.Emitter(b))))
	l.SetMinLevel(Warning)

	ctx := context.Background()
	ctx = alog.AddTags(ctx, "key", "value")
	l.Error(ctx, "I get logged")
	l.Debug(ctx, "I don't get logged")
	const want = `[key=value level=error] I get logged` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got: %#q, want: %#q", got, want)
	}
}
