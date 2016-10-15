package pygmy

import (
	"context"
)

type contextKey int

const (
	shareKey contextKey = iota + 1
)

func setShare(ctx context.Context, obj *share) context.Context {
	return context.WithValue(ctx, shareKey, obj)
}

func getShare(ctx context.Context) *share {
	if obj, ok := ctx.Value(shareKey).(*share); ok {
		return obj
	}
	return nil
}
