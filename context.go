package alog

import (
	"context"
	"encoding/json"
)

type key struct{}

var ctxkey = &key{}

// Tag contains a key and value associated with a log message.  Its value may be
// JSON which may be important to some emitters.
type Tag struct {
	Key    string
	Value  string
	IsJSON bool
}

// AddTags adds paired strings to the set of tags in the Context.
//
// Any unpaired strings are ignored.
func AddTags(ctx context.Context, pairs ...string) context.Context {
	old := fromContext(ctx)
	new := make([]Tag, len(old)+(len(pairs)/2))
	copy(new, old)
	for o := range new[len(old):] {
		new[len(old)+o] = Tag{
			Key:    pairs[o*2],
			Value:  pairs[o*2+1],
			IsJSON: isJSONStructure(pairs[o*2+1]),
		}
	}
	return context.WithValue(ctx, ctxkey, new)
}

// isJSON checks to make sure the input is valid JSON and is either an array or
// object literal.  Note that json.Valid returns true for numeric literals like
// 42.
func isJSONStructure(value string) bool {
	maybeJSON := false
	if len(value) > 0 {
		if value[0] == '{' || value[0] == '[' {
			maybeJSON = true
		}
	}
	return maybeJSON && json.Valid([]byte(value))
}

// fromContext wraps the type assertion coming out of a Context.
func fromContext(ctx context.Context) []Tag {
	if t, ok := ctx.Value(ctxkey).([]Tag); ok {
		return t
	}
	return nil
}
