package otel

import (
	"context"
	"testing"

	"github.com/vimeo/alog/v3/emitter/gkelog"
	"go.opentelemetry.io/otel/api/trace/testtrace"
)

func TestExtractSpanInfo(t *testing.T) {
	outerCtx := context.Background()

	tracer := testtrace.NewTracer()

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
				ctx, unsampledspan := tracer.Start(outerCtx, "foobar")
				sctx := unsampledspan.SpanContext()

				return ctx, gkelog.SpanContext{
					SpanID:  sctx.SpanIDString(),
					TraceID: sctx.TraceIDString(),
					Sampled: false,
				}
			},
		},
		// TODO: Add a sampled variant once it's actually possible to
		// hook up a Sampler to the span creation mechanisms in
		// opentelemetry.
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
