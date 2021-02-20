package log

import (
	"context"
)

// Custom types to avoid key collisions in the context object.
type key int

// Key references an existing log event inside a context.
const Key key = 0

// Inherit tries to get a previously saved log event.
func Inherit(ctx context.Context) Event {
	if ctx == nil {
		return DefaultEvent
	}
	if e, ok := ctx.Value(Key).(Event); ok {
		return e
	}
	return DefaultEvent
}

// Context adds a log event to the current context.
func Context(ctx context.Context, e Event) context.Context {
	ctx = context.WithValue(ctx, Key, e)
	return ctx
}
