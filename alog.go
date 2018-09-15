// Package alog provides a simple logger that has minimal structuring passed via
// context.Context.
package alog // import "github.com/vimeo/alog"

import (
	"context"
	"fmt"
	"os"
	"time"
)

// Default is the the Logger the package-level Print functions use.
var Default = New(To(os.Stderr), WithShortFile(), WithDateFormat(time.RFC3339), WithUTC())

// Print uses the Default Logger to print the supplied string.
func Print(ctx context.Context, v ...interface{}) {
	Default.Output(ctx, 2, fmt.Sprint(v...))
}

// Printf uses the Default Logger to format and then print the supplied string.
func Printf(ctx context.Context, f string, v ...interface{}) {
	Default.Output(ctx, 2, fmt.Sprintf(f, v...))
}
