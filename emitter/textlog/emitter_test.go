package textlog

import (
	"context"
	"os"
	"time"

	"github.com/vimeo/alog/v3"
)

func ExampleEmitter() {
	ctx := context.Background()
	l := alog.New(alog.WithCaller(),
		alog.WithEmitter(Emitter(os.Stdout, WithShortFile(), WithDateFormat(time.RFC3339))),
		alog.OverrideTimestamp(func() time.Time { return time.Time{} }))

	ctx = alog.AddTags(ctx, "allthese", "tags")
	l.Print(ctx, "test")
	// Output:
	// 0001-01-01T00:00:00Z emitter_test.go:18: [allthese=tags] test
}
