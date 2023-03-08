package jsonlog

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/vimeo/alog/v3"
)

var zeroTimeOpt = alog.OverrideTimestamp(func() time.Time { return time.Time{} })

func TestEmitter(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithCaller(), alog.WithEmitter(Emitter(b, WithShortFile())), zeroTimeOpt)

	structuredVal := struct {
		X int `json:"x"`
		Y int `json:"why"`
	}{
		X: 1,
		Y: 42,
	}

	ctx = alog.AddTags(ctx, "allthese", "tags", "andanother", "tag")
	ctx = alog.AddStructuredTags(ctx, alog.STag{Key: "structured", Val: structuredVal}, alog.STag{Key: "other-struct", Val: structuredVal})
	l.Print(ctx, "test")

	want := `{"timestamp":"0001-01-01T00:00:00.000000000Z", "caller":"emitter_test.go:30", "tags":{"allthese":"tags", "andanother":"tag"}, "sTags":{"structured":{"x":1,"why":42}, "other-struct":{"x":1,"why":42}}, "message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if !json.Valid([]byte(got)) {
		t.Errorf("invalid json: %s", got)
	}
}

func TestMessageOnly(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(b, WithDateFormat(""))))

	l.Print(ctx, "test")

	want := `{"message":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestCustomFieldNames(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithCaller(),
		alog.WithEmitter(Emitter(b,
			WithShortFile(),
			WithTimestampField("ts"),
			WithCallerField("called_at"),
			WithMessageField("msg"))), zeroTimeOpt)

	l.Print(ctx, "test")

	want := `{"ts":"0001-01-01T00:00:00.000000000Z", "called_at":"emitter_test.go:66", "msg":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestJSONEscapeValue(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(b, WithDateFormat(""))))

	l.Print(ctx, `"\	`)

	want := `{"message":"\"\\\t"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestJSONEscapeKey(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(b, WithDateFormat(""), WithMessageField("m	s	g"))))

	l.Print(ctx, "test")

	want := `{"m\ts\tg":"test"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestHTMLNoEscape(t *testing.T) {
	b := &bytes.Buffer{}
	ctx := context.Background()
	l := alog.New(alog.WithEmitter(Emitter(b, WithDateFormat(""))))

	l.Print(ctx, "https://vimeo.com")

	want := `{"message":"https://vimeo.com"}` + "\n"
	got := b.String()
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestDuplicateTag(t *testing.T) {
	b := &bytes.Buffer{}
	l := alog.New(alog.WithEmitter(Emitter(b, WithDateFormat(""))))

	// If a caller adds some tags...
	ctx := alog.AddTags(context.Background(), "a", "1", "b", "2")
	// And then adds another tag with the same key...
	ctx = alog.AddTags(ctx, "a", "3")
	// Make sure only the latest one shows up...
	l.Print(ctx, "")
	const want = `{"tags":{"b":"2", "a":"3"}, "message":""}` + "\n"
	if got := b.String(); got != want {
		t.Errorf("got: %#q, want: %#q", got, want)
	}
	// And that it's valid json.
	tgt := struct {
		Message string
		Tags    struct {
			A string
			B string
		}
	}{}
	if err := json.Unmarshal(b.Bytes(), &tgt); err != nil {
		t.Error(err)
	}
}
