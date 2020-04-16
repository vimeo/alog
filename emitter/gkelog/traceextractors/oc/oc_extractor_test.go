package oc

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/vimeo/alog/v3/emitter/gkelog"
	"go.opencensus.io/trace"
)

func TestExtractSpanInfo(t *testing.T) {
	outerCtx := context.Background()
	for _, itbl := range []struct {
		name   string
		ctxGen func() (context.Context, gkelog.SpanContext)
	}{
		{
			name: "nospan",
			ctxGen: func() (context.Context, gkelog.SpanContext) {
				return outerCtx, gkelog.SpanContext{}
			},
		},
		{
			name: "no_sample_span",
			ctxGen: func() (context.Context, gkelog.SpanContext) {
				ctx, unsampledspan := trace.StartSpan(outerCtx, "foobar", trace.WithSampler(trace.NeverSample()))
				sctx := unsampledspan.SpanContext()

				return ctx, gkelog.SpanContext{
					SpanID:  hex.EncodeToString(sctx.SpanID[:]),
					TraceID: hex.EncodeToString(sctx.TraceID[:]),
					Sampled: false,
				}
			},
		},
		{
			name: "sampled_span",
			ctxGen: func() (context.Context, gkelog.SpanContext) {
				ctx, sampledspan := trace.StartSpan(outerCtx, "foobar", trace.WithSampler(trace.AlwaysSample()))
				sctx := sampledspan.SpanContext()

				return ctx, gkelog.SpanContext{
					SpanID:  hex.EncodeToString(sctx.SpanID[:]),
					TraceID: hex.EncodeToString(sctx.TraceID[:]),
					Sampled: true,
				}
			},
		},
	} {

		tbl := itbl
		t.Run(tbl.name, func(t *testing.T) {
			ctx, expectedSctx := tbl.ctxGen()
			sctx := ExtractSpanInfo(ctx)
			if sctx != expectedSctx {
				t.Errorf("expected span context: %+v, got %+v", expectedSctx, sctx)
			}
		})

	}
}
