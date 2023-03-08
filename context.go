package alog

import "context"

type stringKey struct{}
type structuredKey struct{}

var stringCtxKey = stringKey{}
var structuredCtxKey = structuredKey{}

// STag is a structured tag.
type STag struct {
	// A unique key for the structure being logged.
	Key string
	// The structure you would like logged.  This will be marshalled into an
	// appropriate form based on the chosen emitter.
	Val interface{}
}

// AddTags adds paired strings to the set of tags in the Context.
//
// Any unpaired strings are ignored.
func AddTags(ctx context.Context, pairs ...string) context.Context {
	old := tagsFromContext(ctx)
	new := make([][2]string, len(old)+(len(pairs)/2))
	copy(new, old)
	for o := range new[len(old):] {
		new[len(old)+o][0] = pairs[o*2]
		new[len(old)+o][1] = pairs[o*2+1]
	}
	return context.WithValue(ctx, stringCtxKey, new)
}

// tagsFromContext wraps the type assertion coming out of a Context.
func tagsFromContext(ctx context.Context) [][2]string {
	if t, ok := ctx.Value(stringCtxKey).([][2]string); ok {
		return t
	}
	return nil
}

// AddStructuredTags adds tag structures to the Context.
func AddStructuredTags(ctx context.Context, tags ...STag) context.Context {
	oldTags := sTagsFromContext(ctx)

	newTags := append(oldTags[:len(oldTags):len(oldTags)], tags...)

	return context.WithValue(ctx, structuredCtxKey, newTags)
}

// sTagsFromContext wraps the type assertion for structured tags coming out
// of the context.
func sTagsFromContext(ctx context.Context) []STag {
	if tags, ok := ctx.Value(structuredCtxKey).([]STag); ok {
		return tags
	}
	return ([]STag)(nil)
}
