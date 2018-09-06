package alog

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"testing"
)

func Example() {
	ctx := context.Background()
	l := New(To(os.Stdout))

	l.Print(ctx, "test")

	ctx = AddTags(ctx, "more", "context")
	l.Print(ctx, "test")

	ctx = AddTags(ctx, "most", "context")
	l.Print(ctx, "test")
	// Output:
	// test
	// [more=context] test
	// [more=context most=context] test
}

func Example_levels() {
	ctx := context.Background()
	l := New(WithEmitter(func(e *Entry) {
		for _, p := range e.Tags {
			if p[0] != "level" {
				continue
			}
			switch p[1] {
			case "error":
				fmt.Println("ERROR", e.Tags, e.Msg)
				fallthrough
			case "info":
				fallthrough
			case "debug":
				return
			}
		}
	}))
	error := AddTags(ctx, "level", "error")
	info := AddTags(ctx, "level", "info")
	debug := AddTags(ctx, "level", "debug")

	l.Print(debug, "test")
	l.Print(info, "test")
	l.Print(error, "test")
	// Output:
	// ERROR [[level error]] test
}

func ExampleNew() {
	ctx := context.Background()
	l := New(To(os.Stdout), WithPrefix("Example "), WithShortFile())

	ctx = AddTags(ctx, "example", "true")
	l.Print(ctx, "Examples have")
	l.Print(ctx, "weird line numbers")
	// Output:
	// Example alog_test.go:62 [example=true] Examples have
	// Example alog_test.go:63 [example=true] weird line numbers
}

func ExampleWithEmitter() {
	dumper := func(e *Entry) {
		fmt.Printf("%v %s\n", e.Tags, e.Msg)
	}
	ctx := context.Background()
	l := New(WithEmitter(dumper), WithFile())

	ctx = AddTags(ctx, "allthese", "tags")
	l.Print(ctx, "test")
	// Output:
	// [[allthese tags]] test
}

func TestNilOK(t *testing.T) {
	t.Parallel()
	var l *Logger
	ctx := context.Background()

	l.Print(ctx, "this shouldn't explode")
}

func TestIgnoredTag(t *testing.T) {
	t.Parallel()
	buf := &bytes.Buffer{}
	want := "[a=b] test\n"
	l := New(To(buf))

	ctx := AddTags(context.Background(), "a", "b", "unpaired")
	l.Print(ctx, "test")

	if got := buf.String(); want != got {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}

}

func TestToOption(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	buf := &bytes.Buffer{}
	want := "test\n"
	var l *Logger

	l = New(To(buf))
	l.Print(ctx, "test")
	if got := buf.String(); want != got {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}

	buf.Reset()
	l = nil

	l = New(func(l *Logger) {
		WithEmitter(l.EmitText(buf))(l)
	})
	l.Print(ctx, "test")
	if got := buf.String(); want != got {
		t.Fatalf("want: %#q, got: %#q", want, got)
	}
}
