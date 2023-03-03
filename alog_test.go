package alog

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"
)

func Example_levels() {
	ctx := context.Background()
	l := New(WithEmitter(EmitterFunc(func(ctx context.Context, e *Entry) {
		for _, p := range e.Tags {
			if p.Key != "level" {
				continue
			}
			switch p.Value {
			case "error":
				fmt.Println("ERROR", tagsText(e.Tags), e.Msg)
				fallthrough
			case "info":
				fallthrough
			case "debug":
				return
			}
		}
	})))
	error := AddTags(ctx, "level", "error")
	info := AddTags(ctx, "level", "info")
	debug := AddTags(ctx, "level", "debug")

	l.Print(debug, "test")
	l.Print(info, "test")
	l.Print(error, "test")
	// Output:
	// ERROR [[level error]] test
}

func ExampleWithEmitter() {
	dumper := EmitterFunc(func(ctx context.Context, e *Entry) {
		fmt.Printf("%s %s\n", tagsText(e.Tags), e.Msg)
	})
	ctx := context.Background()
	l := New(WithEmitter(dumper))

	ctx = AddTags(ctx, "allthese", "tags")
	l.Print(ctx, "test")
	// Output:
	// [[allthese tags]] test
}

func ExampleWithCaller() {
	dumper := EmitterFunc(func(ctx context.Context, e *Entry) {
		fmt.Printf("%s:%d %s\n", filepath.Base(e.File), e.Line, e.Msg)
	})
	ctx := context.Background()
	l := New(WithEmitter(dumper), WithCaller())

	l.Print(ctx, "test")
	// Output:
	// alog_test.go:61 test
}

func TestOverrideTimestamp(t *testing.T) {
	buf := &bytes.Buffer{}
	dumper := EmitterFunc(func(ctx context.Context, e *Entry) {
		fmt.Fprintf(buf, "%s %s\n", e.Time.Format(time.RFC3339), e.Msg)
	})

	ctx := context.Background()
	l := New(WithEmitter(dumper), OverrideTimestamp(func() time.Time { return time.Time{} }))

	l.Print(ctx, "test")

	want := "0001-01-01T00:00:00Z test\n"
	got := buf.String()
	if got != want {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}
}

func TestNilOK(t *testing.T) {
	t.Parallel()
	var l *Logger
	ctx := context.Background()

	l.Print(ctx, "this shouldn't explode")
}

func tagsText(tags []Tag) string {
	str := "["
	for _, t := range tags {
		str += fmt.Sprintf("[%s %s]", t.Key, t.Value)
	}
	str += "]"
	return str
}

func TestIgnoredTag(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	want := "[[a b]] test\n"
	l := New(WithEmitter(EmitterFunc(func(ctx context.Context, e *Entry) {
		fmt.Fprintf(buf, "%s %s\n", tagsText(e.Tags), e.Msg)
	})))

	ctx := AddTags(context.Background(), "a", "b", "unpaired")
	l.Print(ctx, "test")

	if got := buf.String(); want != got {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}
}

func TestStdLogger(t *testing.T) {
	buf := &bytes.Buffer{}
	dumper := EmitterFunc(func(ctx context.Context, e *Entry) {
		fmt.Fprintf(buf, "%s %s %s\n", e.Time.Format(time.RFC3339), tagsText(e.Tags), e.Msg)
	})

	l := New(WithEmitter(dumper), OverrideTimestamp(func() time.Time { return time.Time{} }))

	ctx := AddTags(context.Background(), "a", "b")
	sl := l.StdLogger(ctx)
	sl.Printf("test")

	want := "0001-01-01T00:00:00Z [[a b]] test\n"
	got := buf.String()
	if got != want {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}
}
