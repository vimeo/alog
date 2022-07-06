module github.com/vimeo/alog/emitter/gkelog/traceextractors/otel

go 1.18

require (
	github.com/vimeo/alog/v3 v3.5.0
	go.opentelemetry.io/otel v0.2.2
)

require google.golang.org/grpc v1.47.0 // indirect

replace github.com/vimeo/alog/v3 => ../../../../
