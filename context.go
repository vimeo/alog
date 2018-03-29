package alog

import "context"

type key struct{}

var ctxkey = &key{}

// AddTags adds paired strings to the set of tags in the Context.
//
// Any unpaired strings are ignored.
func AddTags(ctx context.Context, pairs ...string) context.Context {
	old := fromContext(ctx)
	new := make([][2]string, len(old)+(len(pairs)/2))
	copy(new, old)
	for o := range new[len(old):] {
		new[len(old)+o][0] = pairs[o*2]
		new[len(old)+o][1] = pairs[o*2+1]
	}
	return context.WithValue(ctx, ctxkey, new)
}

// fromContext wraps the type assertion coming out of a Context.
func fromContext(ctx context.Context) [][2]string {
	if t, ok := ctx.Value(ctxkey).([][2]string); ok {
		return t
	}
	return nil
}
