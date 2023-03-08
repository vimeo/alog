package testlog

import (
	"context"
	"testing"

	"github.com/vimeo/alog/v3"
)

func TestEmitter(t *testing.T) {

	ctx := context.Background()
	l := DefaultLogger(t, WithShortFile())

	structured := struct {
		X int
	}{
		X: 1,
	}

	ctx = alog.AddTags(ctx, "test", "tags")
	ctx = alog.AddStructuredTags(ctx, alog.STag{Key: "structured", Val: structured})
	l.Print(ctx, "testMessage")

	// Output
	// logger.go:40: emitter_test.go testMessage [[test tags]] [{Key:structured Val:{X:1}}]
}
