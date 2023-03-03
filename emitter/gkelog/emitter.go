package gkelog

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/vimeo/alog/v3"
	"github.com/vimeo/alog/v3/emitter/internal"
)

// DefaultLogger is a *alog.Logger that logs to stderr
var DefaultLogger = alog.New(alog.WithEmitter(Emitter()))

type contextKey string

var (
	severityKey = contextKey("severity")
	requestKey  = contextKey("request")
	statusKey   = contextKey("status")
	latencyKey  = contextKey("latency")
	traceKey    = contextKey("trace")
	spanKey     = contextKey("span")
)

// WithSeverity returns a copy of parent with the specified severity value.
func WithSeverity(parent context.Context, severity string) context.Context {
	return context.WithValue(parent, severityKey, severity)
}

// WithRequest returns a copy of parent with the specified http.Request value.
// It also calls WithRequestTrace to add trace information to the context.
func WithRequest(parent context.Context, req *http.Request) context.Context {
	ctx := context.WithValue(parent, requestKey, req)
	ctx = WithRequestTrace(ctx, req)
	return ctx
}

// WithTrace returns a copy of parent with the specified Trace ID value.
func WithTrace(parent context.Context, trace string) context.Context {
	return context.WithValue(parent, traceKey, trace)
}

// WithSpan returns a copy of parent with the specified Span ID value.
//
// This should be a 8-byte hex string (16 digits). Note that some
// Google services like load balancers use a 64-bit decimal number instead
// of hexidecimal. Those values can be converted to the correct format
// with SpanDecimalToHex.
func WithSpan(parent context.Context, span string) context.Context {
	return context.WithValue(parent, spanKey, span)
}

// SpanDecimalToHex converts a decimal Span ID value to hexidecimal.
func SpanDecimalToHex(spanID uint64) string {
	spanIDHex := strconv.FormatUint(spanID, 16)
	if len(spanIDHex) < 16 {
		spanIDHex = strings.Repeat("0", 16-len(spanIDHex)) + spanIDHex
	}
	return spanIDHex
}

// TraceFromRequest returns a trace and/or span from a http.Request.
func TraceFromRequest(req *http.Request) (trace string, span string) {
	traceHeader := req.Header.Get("X-Cloud-Trace-Context")
	if traceHeader != "" {
		traceSpan := strings.Split(traceHeader, "/")

		trace = traceSpan[0]
		if len(traceSpan) > 1 {
			spanID, err := strconv.ParseUint(traceSpan[1], 10, 64)
			if err == nil {
				span = SpanDecimalToHex(spanID)
			}
		}
	}
	return
}

// WithRequestTrace returns a copy of parent with the trace information from
// the specified http.Request.
func WithRequestTrace(parent context.Context, req *http.Request) context.Context {
	ctx := parent
	trace, span := TraceFromRequest(req)
	if trace != "" {
		ctx = WithTrace(ctx, trace)
	}
	if span != "" {
		ctx = WithSpan(ctx, span)
	}
	return ctx
}

// WithRequestStatus returns a copy of the parent with the specified HTTP return status code.
func WithRequestStatus(parent context.Context, status int) context.Context {
	return context.WithValue(parent, statusKey, status)
}

// WithRequestLatency returns a copy of the parent with the specified HTTP request latency value.
func WithRequestLatency(parent context.Context, latency time.Duration) context.Context {
	return context.WithValue(parent, latencyKey, latency)
}

func jsonString(w *bytes.Buffer, s string) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.Encode(s)
	w.Truncate(w.Len() - 1)
}

func jsonKey(w *bytes.Buffer, s string) {
	jsonString(w, s)
	w.WriteByte(':')
}

var reservedKeys = map[string]bool{
	"httpHeaders":                           true,
	"httpQuery":                             true,
	"httpRequest":                           true,
	"logging.googleapis.com/sourceLocation": true,
	"logging.googleapis.com/spanId":         true,
	"logging.googleapis.com/trace":          true,
	"message":                               true,
	"severity":                              true,
	"time":                                  true,
}

func defaultTraceExtractor(ctx context.Context) SpanContext {
	sctx := SpanContext{}

	traceV := ctx.Value(traceKey)
	if traceV != nil {
		sctx.TraceID = traceV.(string)
		sctx.Sampled = true
	}
	spanV := ctx.Value(spanKey)
	if spanV != nil {
		sctx.SpanID = spanV.(string)
		sctx.Sampled = true
	}
	return sctx
}

func jsonTrace(ctx context.Context, o *Options, w *bytes.Buffer) {
	if o.spanExtractor == nil {
		return
	}
	sctx := o.spanExtractor(ctx)
	if sctx.TraceID != "" {
		jsonKey(w, "logging.googleapis.com/trace")
		jsonString(w, sctx.TraceID)
		w.WriteString(", ")
	}

	if sctx.SpanID != "" {
		jsonKey(w, "logging.googleapis.com/spanId")
		jsonString(w, sctx.SpanID)
		w.WriteString(", ")
	}
	if sctx.TraceID != "" || sctx.SpanID != "" {
		jsonKey(w, "logging.googleapis.com/trace_sampled")
		w.WriteString(strconv.FormatBool(sctx.Sampled))
		w.WriteString(", ")
	}
}

func sortedMapKeys(m map[string][]string) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

var skipHeaders = map[string]bool{
	"Referer":               true,
	"Referrer":              true,
	"User-Agent":            true,
	"X-Cloud-Trace-Context": true,
}

func jsonHTTPRequest(ctx context.Context, w *bytes.Buffer) {
	var (
		request *http.Request
		status  int
		latency time.Duration
	)

	reqV := ctx.Value(requestKey)
	if reqV != nil {
		request = reqV.(*http.Request)
	}
	statusV := ctx.Value(statusKey)
	if statusV != nil {
		status = statusV.(int)
	}
	latencyV := ctx.Value(latencyKey)
	if latencyV != nil {
		latency = latencyV.(time.Duration)
	}

	if request == nil && status <= 0 && latency == 0 {
		return
	}

	jsonKey(w, "httpRequest")
	w.WriteByte('{')

	if status > 0 {
		jsonKey(w, "status")
		internal.Itoa(w, uint(status))
		if latency > 0 || request != nil {
			w.WriteString(", ")
		}
	}

	if latency > 0 {
		jsonKey(w, "latency")
		latencyStr := strconv.FormatFloat(latency.Seconds(), 'f', -1, 64) + "s"
		jsonString(w, latencyStr)
		if request != nil {
			w.WriteString(", ")
		}
	}

	if request != nil {
		jsonKey(w, "requestMethod")
		jsonString(w, request.Method)
		w.WriteString(", ")

		u := *request.URL
		u.Fragment = ""
		jsonKey(w, "requestUrl")
		jsonString(w, u.String())

		if request.UserAgent() != "" {
			w.WriteString(", ")
			jsonKey(w, "userAgent")
			jsonString(w, request.UserAgent())
		}

		if request.Referer() != "" {
			w.WriteString(", ")
			jsonKey(w, "referer")
			jsonString(w, request.Referer())
		}
	}

	w.WriteByte('}')
	w.WriteString(", ")

	if request == nil {
		return
	}

	headerKeys := sortedMapKeys(request.Header)
	if len(headerKeys) > 0 {
		jsonKey(w, "httpHeaders")
		w.WriteByte('{')
		i := 0
		for _, h := range headerKeys {
			if skipHeaders[h] {
				continue
			}
			if i > 0 {
				w.WriteString(", ")
			}
			jsonKey(w, h)
			v := request.Header[h]
			w.WriteByte('[')
			for j, v0 := range v {
				if j > 0 {
					w.WriteString(", ")
				}
				jsonString(w, v0)
			}
			w.WriteByte(']')
			i++
		}
		w.WriteByte('}')
		w.WriteString(", ")
	}

	query := request.URL.Query()
	queryKeys := sortedMapKeys(query)
	if len(queryKeys) > 0 {
		jsonKey(w, "httpQuery")
		w.WriteByte('{')
		for i, q := range queryKeys {
			if i > 0 {
				w.WriteString(", ")
			}
			jsonKey(w, q)
			v := query[q]
			w.WriteByte('[')
			for j, v0 := range v {
				if j > 0 {
					w.WriteString(", ")
				}
				jsonString(w, v0)
			}
			w.WriteByte(']')
		}
		w.WriteByte('}')
		w.WriteString(", ")
	}

}

// Emitter emits log messages as single lines of JSON.
//
// Logs are output to w. Every entry generates a single Write call to w, and
// calls are serialized.
func Emitter(opt ...Option) alog.Emitter {
	o := &Options{
		spanExtractor: defaultTraceExtractor,
	}
	for _, option := range opt {
		option(o)
	}
	if o.reqWriter == nil {
		o.reqWriter = os.Stdout
	}
	if o.appWriter == nil {
		o.appWriter = os.Stderr
	}

	wReq := internal.NewSerializedWriter(o.reqWriter)
	wApp := internal.NewSerializedWriter(o.appWriter)

	return alog.EmitterFunc(func(ctx context.Context, e *alog.Entry) {
		b := internal.GetBuffer()
		defer internal.PutBuffer(b)

		b.WriteByte('{')

		jsonKey(b, "time")
		jsonString(b, e.Time.UTC().Format(time.RFC3339Nano))
		b.WriteString(", ")

		severity := ctx.Value(severityKey)
		if severity != nil {
			jsonKey(b, "severity")
			jsonString(b, severity.(string))
			b.WriteString(", ")
		}

		jsonHTTPRequest(ctx, b)

		jsonTrace(ctx, o, b)

		tagClean := make(map[string]int, len(e.Tags))
		for i, tag := range e.Tags {
			tagClean[tag.Key] = i
		}
		for i, tag := range e.Tags {
			if tagClean[tag.Key] != i || reservedKeys[tag.Key] {
				continue
			}
			jsonKey(b, tag.Key)
			if tag.IsJSON {
				b.WriteString(tag.Value)
			} else {
				jsonString(b, tag.Value)
			}
			b.WriteString(", ")
		}

		if e.File != "" {
			jsonKey(b, "logging.googleapis.com/sourceLocation")
			b.WriteByte('{')
			jsonKey(b, "file")
			f := e.File
			if o.shortfile {
				f = path.Base(f)
			}
			jsonString(b, f)
			b.WriteString(", ")
			jsonKey(b, "line")
			b.WriteByte('"')
			internal.Itoa(b, uint(e.Line))
			b.WriteByte('"')
			b.WriteByte('}')
			b.WriteString(", ")
		}

		jsonKey(b, "message")
		jsonString(b, e.Msg)

		b.WriteString("}\n")

		if ctx.Value(requestKey) == nil {
			wApp.Write(b.Bytes())
		} else {
			wReq.Write(b.Bytes())
		}
	})
}
