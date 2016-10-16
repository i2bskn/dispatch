package pygmy

import (
	"context"
)

type contextKey int

const (
	shareKey contextKey = iota + 1
	paramKey
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

func setParam(ctx context.Context, param map[string]string) context.Context {
	return context.WithValue(ctx, paramKey, param)
}

func Param(ctx context.Context, key string) string {
	if params, ok := ctx.Value(paramKey).(map[string]string); ok {
		return params[key]
	}
	return ""
}
