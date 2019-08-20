package leveled

import (
	"context"
	"fmt"

	"github.com/vimeo/alog/v3"
)

// Logger is an interface that implements logging functions for different levels of severity.
type Logger interface {
	// Debug logs debugging or trace information.
	Debug(ctx context.Context, f string, v ...interface{})

	// Info logs normal information.
	Info(ctx context.Context, f string, v ...interface{})

	// Warning logs a message that indicates a potential problem.
	Warning(ctx context.Context, f string, v ...interface{})

	// Error logs an error that indicates a likely problem.
	Error(ctx context.Context, f string, v ...interface{})

	// Critical logs an error that indicates a definite problem.
	Critical(ctx context.Context, f string, v ...interface{})
}

type defaultLogger struct {
	*alog.Logger
}

// Default returns a Logger that wraps the provided `alog.Logger`.
// It adds a generic "level" tag that is not intended for special use by
// emitters. If emitter-specific functionality around log levels is needed,
// a different implementation of Logger should be used.
func Default(logger *alog.Logger) Logger {
	return &defaultLogger{
		Logger: logger,
	}
}

// Debug implements Logger.Debug
func (d *defaultLogger) Debug(ctx context.Context, f string, v ...interface{}) {
	alog.AddTags(ctx, "level", "debug")
	d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
}

// Info implements Logger.Info
func (d *defaultLogger) Info(ctx context.Context, f string, v ...interface{}) {
	alog.AddTags(ctx, "level", "info")
	d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
}

// Warning implements Logger.Warning
func (d *defaultLogger) Warning(ctx context.Context, f string, v ...interface{}) {
	alog.AddTags(ctx, "level", "warning")
	d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
}

// Error implements Logger.Error
func (d *defaultLogger) Error(ctx context.Context, f string, v ...interface{}) {
	alog.AddTags(ctx, "level", "error")
	d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
}

// Critical implements Logger.Critical
func (d *defaultLogger) Critical(ctx context.Context, f string, v ...interface{}) {
	alog.AddTags(ctx, "level", "critical")
	d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
}
