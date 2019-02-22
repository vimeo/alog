package alog

import "context"

type tagsKey struct{}

var tagsCtxKey = &tagsKey{}

type loggerKey struct{}

var loggerCtxKey = &loggerKey{}

// AddTags adds paired strings to the set of tags in the Context.
//
// Any unpaired strings are ignored.
func AddTags(ctx context.Context, pairs ...string) context.Context {
	old := tagsFromContext(ctx)
	new := make([][2]string, len(old)+(len(pairs)/2))
	copy(new, old)
	for o := range new[len(old):] {
		new[len(old)+o][0] = pairs[o*2]
		new[len(old)+o][1] = pairs[o*2+1]
	}
	return context.WithValue(ctx, tagsCtxKey, new)
}

// tagsFromContext wraps the type assertion for tags coming out of a Context.
func tagsFromContext(ctx context.Context) [][2]string {
	if t, ok := ctx.Value(tagsCtxKey).([][2]string); ok {
		return t
	}
	return nil
}

// AddLogger adds a logger instance to the passed Context.
func AddLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerCtxKey, l)
}

// LoggerFromContext wraps the type assertion for the logger coming out of a Context.
func LoggerFromContext(ctx context.Context) *Logger {
	if l, ok := ctx.Value(loggerCtxKey).(*Logger); ok {
		return l
	}

	return nil
}
