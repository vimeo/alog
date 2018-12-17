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

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/sourceLocation":{"file":"emitter_test.go", "line":"23"}, "message":"test"}` + "\n"
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

	ctx = alog.AddTags(ctx, "allthese", "tags", "andanother", "tag")
	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "allthese":"tags", "andanother":"tag", "message":"test"}` + "\n"
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

func TestRequest(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(WithWriter(b))), zeroTimeOpt)

	req := httptest.NewRequest(http.MethodGet, "/test/endpoint", strings.NewReader("this is a test"))
	req.Header.Set("User-Agent", "curl/7.54.0")
	req.Header.Set("Referer", "https://vimeo.com")
	req.Header.Set("X-Cloud-Trace-Context", "a2fbf27a2ed90077e0d4af0e40a241f9/12690385211238481741")
	ctx = WithRequest(ctx, req)

	l.Print(ctx, "test")

	want := `{"time":"0001-01-01T00:00:00Z", "httpRequest":{"requestMethod":"GET", "requestUrl":"/test/endpoint", "userAgent":"curl/7.54.0", "referer":"https://vimeo.com"}, "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "message":"test"}` + "\n"
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

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "message":"test"}` + "\n"
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

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "message":"test"}` + "\n"
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

	want := `{"time":"0001-01-01T00:00:00Z", "logging.googleapis.com/trace":"a2fbf27a2ed90077e0d4af0e40a241f9", "logging.googleapis.com/spanId":"b01d4e1cf2bd7f4d", "message":"test"}` + "\n"
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
