package leveled

import (
	"context"
	"fmt"

	"github.com/vimeo/alog/v3"
)

// Level represents the severity of a log message.
type Level uint8

const (
	Debug    Level = iota // debug
	Info                  // info
	Warning               // warning
	Error                 // error
	Critical              // critical
)

// LevelKey is the tag key associated with a level.
const LevelKey = "level"

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

// FilteredLogger amends the Logger interface with a Log method that accepts the
// level to log at.
type FilteredLogger interface {
	Logger
	Log(ctx context.Context, level Level, f string, v ...interface{})
	SetMinLevel(level Level)
}

type defaultLogger struct {
	*alog.Logger

	// Indicates the minimum level to log at.  If MinLevel is greater than the
	// level of a given log message, the log message will be suppressed.
	MinLevel Level
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

// Filtered returns a logger that allows for setting the minimum level.
func Filtered(logger *alog.Logger) FilteredLogger {
	return &defaultLogger{
		Logger: logger,
	}
}

// Debug implements Logger.Debug
func (d *defaultLogger) Debug(ctx context.Context, f string, v ...interface{}) {
	d.Log(ctx, Debug, f, v...)
}

// Info implements Logger.Info
func (d *defaultLogger) Info(ctx context.Context, f string, v ...interface{}) {
	d.Log(ctx, Info, f, v...)
}

// Warning implements Logger.Warning
func (d *defaultLogger) Warning(ctx context.Context, f string, v ...interface{}) {
	d.Log(ctx, Warning, f, v...)
}

// Error implements Logger.Error
func (d *defaultLogger) Error(ctx context.Context, f string, v ...interface{}) {
	d.Log(ctx, Error, f, v...)
}

// Critical implements Logger.Critical
func (d *defaultLogger) Critical(ctx context.Context, f string, v ...interface{}) {
	d.Log(ctx, Critical, f, v...)
}

// Log implements FilteredLogger.Log
func (d *defaultLogger) Log(ctx context.Context, level Level, f string, v ...interface{}) {
	if level >= d.MinLevel {
		ctx = alog.AddTags(ctx, LevelKey, level.String())
		d.Logger.Output(ctx, 3, fmt.Sprintf(f, v...))
	}
}

// SetMinLevel sets the minimum level that will be logged and implements FilteredLogger.
func (d *defaultLogger) SetMinLevel(level Level) {
	d.MinLevel = level
}

//go:generate go run golang.org/x/tools/cmd/stringer@latest -type Level -linecomment
