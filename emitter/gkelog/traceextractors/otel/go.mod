module github.com/vimeo/alog/emitter/gkelog/traceextractors/otel

go 1.12

require (
	github.com/vimeo/alog/v3 v3.5.0
	go.opentelemetry.io/otel v0.2.2
	google.golang.org/grpc v1.56.3 // indirect
)

replace github.com/vimeo/alog/v3 => ../../../../
