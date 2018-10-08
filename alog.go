// Package alog provides a simple logger that has minimal structuring passed via
// context.Context.
package alog

import (
	"context"
	"fmt"
	"time"
)

var defaultEmitter = EmitterFunc(func(ctx context.Context, e *Entry) {
	tags := ""
	if len(e.Tags) > 0 {
		tags = " " + fmt.Sprintf("%v", e.Tags)
	}
	fmt.Printf("%s%s %s\n", e.Time.Format(time.RFC3339), tags, e.Msg)
})

// Default is the the Logger the package-level Print functions use.
var Default = New(WithEmitter(defaultEmitter))

// Print uses the Default Logger to print the supplied string.
func Print(ctx context.Context, v ...interface{}) {
	Default.Output(ctx, 2, fmt.Sprint(v...))
}

// Printf uses the Default Logger to format and then print the supplied string.
func Printf(ctx context.Context, f string, v ...interface{}) {
	Default.Output(ctx, 2, fmt.Sprintf(f, v...))
}
