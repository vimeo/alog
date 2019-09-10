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
