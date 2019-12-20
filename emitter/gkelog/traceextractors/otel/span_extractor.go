package otel

import (
	"context"

	"github.com/vimeo/alog/v3/emitter/gkelog"

	"go.opentelemetry.io/otel/api/core"
	"go.opentelemetry.io/otel/api/trace"
)

// ExtractSpanInfo extracts a span ID from the passed context if there is an
// opentelemetry trace-span embedded within.
// Returns spanID, traceID, isSampled
func ExtractSpanInfo(ctx context.Context) gkelog.SpanContext {
	span := trace.SpanFromContext(ctx)
	if span == nil {
		return gkelog.SpanContext{}
	}
	sctx := span.SpanContext()
	if sctx == core.EmptySpanContext() {
		return gkelog.SpanContext{}
	}
	return gkelog.SpanContext{
		SpanID:  sctx.SpanIDString(),
		TraceID: sctx.TraceIDString(),
		Sampled: sctx.IsSampled(),
	}

}

// Provide a guarantee that ExtractSpanInfo matches the interface definition for gkelog.TraceSpanExtractor
var _ gkelog.TraceSpanExtractor = ExtractSpanInfo
