package gkelog

import (
	"context"
	"io"
)

// Options holds option values.
type Options struct {
	reqWriter     io.Writer
	appWriter     io.Writer
	spanExtractor TraceSpanExtractor
	shortfile     bool
}

// Option sets an option for the emitter.
//
// Options are applied in the order specified.
type Option func(*Options)

// WithWriter sets the writer to use for log output.
func WithWriter(w io.Writer) Option {
	return func(o *Options) {
		o.reqWriter = w
		o.appWriter = w
	}
}

// WithWriters sets the writers to use for log output.
//
// req is used for log entries that have a Request.
// app is used for all other log entries.
// In GKE the typical usage would be WithWriters(os.Stdout, os.Stderr).
func WithWriters(req io.Writer, app io.Writer) Option {
	return func(o *Options) {
		o.reqWriter = req
		o.appWriter = app
	}
}

// WithShortFile indicates to only use the filename instead of the full file
// path for sourceLocation.
func WithShortFile() Option {
	return func(o *Options) { o.shortfile = true }
}

// SpanContext contains the necessary trace-context data to populate the
// logging.googleapis.com/spanId, logging.googleapis.com/trace and
// logging.googleapis.com/trace_sampled fields.
// See https://cloud.google.com/logging/docs/agent/configuration#special-fields
// for details on these fields.
type SpanContext struct {
	SpanID  string
	TraceID string
	Sampled bool
}

// TraceSpanExtractor implementations extract a trace spanID from the passed
// context.
// The return values should be: spanID, traceID, isSampled
type TraceSpanExtractor func(context.Context) SpanContext

// WithTraceSpanExtractor registers a trace-span extractor so trace-span IDs
// from the context (as found by the extractor) are placed in the appropriate
// fields to be correlated by stackdriver tracing.
func WithTraceSpanExtractor(extractor TraceSpanExtractor) Option {
	return func(o *Options) { o.spanExtractor = extractor }
}
