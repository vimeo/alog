package gkelog

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/vimeo/alog/v3"
)

var zeroTimeOpt = alog.OverrideTimestamp(func() time.Time { return time.Time{} })

// Keep this function at the top of the file so that the line number doesn't change too often
func TestCaller(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithCaller(), alog.WithEmitter(Emitter(WithWriter(b), WithShortFile())), zeroTimeOpt)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/sourceLocation":{"file":"emitter_test.go", "line":"24"}, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestEmitter(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLabels(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	structured := alog.STag{
		Key: "structured",
		Val: struct {
			X int `json:"x"`
		}{
			X: 1,
		},
	}

	structured2 := alog.STag{
		Key: "structured-again",
		Val: struct {
			Y string `json:"val"`
		}{
			Y: "foo",
		},
	}

	ctx = alog.AddTags(ctx, "allthese", "tags", "andanother", "tag")
	ctx = alog.AddStructuredTags(ctx, structured, structured2)
	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "allthese":"tags", "andanother":"tag", "structured":{"x":1}, "structured-again":{"val":"foo"}, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestSeverity(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = WithSeverity(ctx, SeverityError)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "severity":"ERROR", "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLogSeverity(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	LogInfo(ctx, l, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "severity":"INFO", "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestMinLogSeverity(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := WithMinSeverity(context.Background(), SeverityWarning)
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	LogInfo(ctx, l, "NOT LOGGED") // because Info is lower than Warning
	LogError(ctx, l, "LOGGED")

	want := `{"time":"0001-01-01T00:00:00Z", "severity":"ERROR", "message":"LOGGED"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRequest(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	req := httptest.NewRequest(http.MethodGet, "/test/endpoint?q=1&c=pink&c=red", strings.NewReader("this is a test"))
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.Header.Set("Referer", "https://vimeo.com")
	req.Header.Set("X-Cloud-Trace-Context", "a2fbf27a2ed90077e0d4af0e40a241f9/12690385211238481741")
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Dnt", "1")
	req.Header.Set("X-Varnish", "731698977")
	req.Header.Add("X-Varnish", "4193052513")
	ctx = WithRequest(ctx, req)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "httpRequest":{"requestMethod":"GET", "requestUrl":"/test/endpoint?q=1&c=pink&c=red", "userAgent":"curl/7.54.0", "referer":"https://vimeo.com"}, "httpHeaders":{"Content-Type":["text/plain"], "Dnt":["1"], "X-Varnish":["731698977", "4193052513"]}, "httpQuery":{"c":["pink", "red"], "q":["1"]}, "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "logging.googleapis.com/trace_sampled":true, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestTrace(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = WithTrace(ctx, "a2fbf27a2ed90077e0d4af0e40a241f9")

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "logging.googleapis.com/trace_sampled":true, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestSpanDecimalToHex(t *testing.T) {
	data := []string{
		"7f1935142c348765",
		"0eb8507a8410c0ec",
		"000000008410c0ec",
		"0000000000000000",
	}
	conv := []string{}
	for _, hex := range data {
		dec, _ := strconv.ParseUint(hex, 16, 64)
		conv = append(conv, SpanDecimalToHex(dec))
	}

	want := strings.Join(data, " ")
	got := strings.Join(conv, " ")
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestSpan(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = WithSpan(ctx, SpanDecimalToHex(12690385211238481741))

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "logging.googleapis.com/trace_sampled":true, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestRequestTrace(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	req := httptest.NewRequest(http.MethodGet, "/test/endpoint", strings.NewReader("this is a test"))
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.Header.Set("Referer", "https://vimeo.com")
	req.Header.Set("X-Cloud-Trace-Context", "a2fbf27a2ed90077e0d4af0e40a241f9/12690385211238481741")
	ctx = WithRequestTrace(ctx, req)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "logging.googleapis.com/trace_sampled":true, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestStatus(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = WithRequestStatus(ctx, http.StatusBadRequest)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "httpRequest":{"status":400}, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestLatency(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = WithRequestLatency(ctx, 1549284472*time.Nanosecond)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "httpRequest":{"latency":"1.549284472s"}, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestWriters(t *testing.T) {
	b0 := &bytes.Buffer{}
	b1 := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriters(b0, b1))), zeroTimeOpt)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rctx := WithRequestStatus(WithRequest(ctx, req), http.StatusOK)

	b0.WriteString("b0: ")
	l.Print(rctx, "test")
	b1.WriteString("b1: ")
	l.Print(ctx, "test")

	want := `b0: {"time":"0001-01-01T00:00:00Z", "httpRequest":{"status":200, "requestMethod":"GET", "requestUrl":"/test"}, "message":"test"}` + "\n" + `b1: {"time":"0001-01-01T00:00:00Z", "message":"test"}` + "\n"
	got := b0.String() + b1.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestDuplicateTags(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = alog.AddTags(ctx, "a", "v1", "a", "v2")
	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "a":"v2", "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestReservedKeys(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	ctx = alog.AddTags(ctx, "time", "2019-08-21T19:02:23Z")
	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}
