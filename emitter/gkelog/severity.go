package gkelog

import (
	"context"
	"fmt"

	"github.com/vimeo/alog/v3"
	"github.com/vimeo/alog/v3/leveled"
)

// Severity levels
const (
	SeverityDefault   = "DEFAULT"   // The log entry has no assigned severity level.
	SeverityDebug     = "DEBUG"     // Debug or trace information.
	SeverityInfo      = "INFO"      // Routine information, such as ongoing status or performance.
	SeverityNotice    = "NOTICE"    // Normal but significant events, such as start up, shut down, or a configuration change.
	SeverityWarning   = "WARNING"   // Warning events might cause problems.
	SeverityError     = "ERROR"     // Error events are likely to cause problems.
	SeverityCritical  = "CRITICAL"  // Critical events cause more severe problems or outages.
	SeverityAlert     = "ALERT"     // A person must take an action immediately.
	SeverityEmergency = "EMERGENCY" // One or more systems are unusable.
)

// the priority of each severity level
var severityPriority map[string]uint8 = map[string]uint8{
	SeverityDebug:     0,
	SeverityInfo:      1,
	SeverityNotice:    2,
	SeverityWarning:   3,
	SeverityError:     4,
	SeverityCritical:  5,
	SeverityAlert:     6,
	SeverityEmergency: 7,
	SeverityDefault:   8, // default will almost always log and should probably not be used
}

// Separate private function so that LogSeverity and the other logs functions
// will have the same stack frame depth and thus use the same calldepth value.
// See https://golang.org/pkg/runtime/#Caller and
// https://godoc.org/github.com/vimeo/alog#Logger.Output
func logSeverity(ctx context.Context, logger *alog.Logger, s string, f string, v ...interface{}) {
	minSeverity := uint8(0)
	if minSeverityVal := ctx.Value(minSeverityKey); minSeverityVal != nil {
		minSeverity = severityPriority[minSeverityVal.(string)]
	}
	if severityPriority[s] >= minSeverity {
		ctx = WithSeverity(ctx, s)
		logger.Output(ctx, 3, fmt.Sprintf(f, v...))
	}
}

// LogSeverity writes a log entry using the specified severity
func LogSeverity(ctx context.Context, logger *alog.Logger, severity string, f string, v ...interface{}) {
	logSeverity(ctx, logger, severity, f, v...)
}

// LogDebug writes a log entry using SeverityDebug
func LogDebug(ctx context.Context, logger *alog.Logger, f string, v ...interface{}) {
	logSeverity(ctx, logger, SeverityDebug, f, v...)
}

// LogInfo writes a log entry using SeverityInfo
func LogInfo(ctx context.Context, logger *alog.Logger, f string, v ...interface{}) {
	logSeverity(ctx, logger, SeverityInfo, f, v...)
}

// LogWarning writes a log entry using SeverityWarning
func LogWarning(ctx context.Context, logger *alog.Logger, f string, v ...interface{}) {
	logSeverity(ctx, logger, SeverityWarning, f, v...)
}

// LogError writes a log entry using SeverityError
func LogError(ctx context.Context, logger *alog.Logger, f string, v ...interface{}) {
	logSeverity(ctx, logger, SeverityError, f, v...)
}

// LogCritical writes a log entry using SeverityCritical
func LogCritical(ctx context.Context, logger *alog.Logger, f string, v ...interface{}) {
	logSeverity(ctx, logger, SeverityCritical, f, v...)
}

type severityLogger struct {
	*alog.Logger
}

func NewSeverityLogger(logger *alog.Logger) leveled.Logger {
	return &severityLogger{
		Logger: logger,
	}
}

func (sl *severityLogger) Debug(ctx context.Context, f string, v ...interface{}) {
	logSeverity(ctx, sl.Logger, SeverityDebug, f, v...)
}

func (sl *severityLogger) Info(ctx context.Context, f string, v ...interface{}) {
	logSeverity(ctx, sl.Logger, SeverityInfo, f, v...)
}

func (sl *severityLogger) Warning(ctx context.Context, f string, v ...interface{}) {
	logSeverity(ctx, sl.Logger, SeverityWarning, f, v...)
}

func (sl *severityLogger) Error(ctx context.Context, f string, v ...interface{}) {
	logSeverity(ctx, sl.Logger, SeverityError, f, v...)
}

func (sl *severityLogger) Critical(ctx context.Context, f string, v ...interface{}) {
	logSeverity(ctx, sl.Logger, SeverityCritical, f, v...)
}
