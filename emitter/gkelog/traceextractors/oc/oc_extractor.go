package oc

import (
	"context"
	"encoding/hex"

	"github.com/vimeo/alog/v3/emitter/gkelog"

	"go.opencensus.io/trace"
)

// ExtractSpanInfo extracts a span ID from the passed context if there is an
// opencensus trace-span embedded within.
// Returns spanID, traceID, isSampled
func ExtractSpanInfo(ctx context.Context) gkelog.SpanContext {
	span := trace.FromContext(ctx)
	if span == nil {
		return gkelog.SpanContext{}
	}
	sctx := span.SpanContext()

	return gkelog.SpanContext{
		SpanID:  hex.EncodeToString(sctx.SpanID[:]),
		TraceID: hex.EncodeToString(sctx.TraceID[:]),
		Sampled: sctx.IsSampled(),
	}

}

// Provide a guarantee that ExtractSpanInfo matches the interface definition for gkelog.TraceSpanExtractor
var _ gkelog.TraceSpanExtractor = ExtractSpanInfo
