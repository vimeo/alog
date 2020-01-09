package testlog

import (
	"context"
	"testing"

	"github.com/vimeo/alog/v3"
)

func TestEmitter(t *testing.T) {

	ctx := context.Background()
	l := DefaultLogger(t, WithShortFile())

	ctx = alog.AddTags(ctx, "test", "tags")
	l.Print(ctx, "testMessage")

	// Output
	// logger.go:40: emitter_test.go testMessage [[test tags]]
}
