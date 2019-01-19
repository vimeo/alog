package gkelog

import (
	"context"
	"fmt"

	"github.com/vimeo/alog/v3"
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

// Separate private function so that LogSeverity and the other logs functions
// will have the same stack frame depth and thus use the same calldepth value.
// See https://golang.org/pkg/runtime/#Caller and
// https://godoc.org/github.com/vimeo/alog#Logger.Output
func logSeverity(ctx context.Context, logger *alog.Logger, s string, f string, v ...interface{}) {
	ctx = WithSeverity(ctx, s)
	logger.Output(ctx, 3, fmt.Sprintf(f, v...))
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
