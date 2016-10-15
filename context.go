package pygmy

import (
	"context"
	"net/http"
)

type contextKey int

const (
	handlerKey contextKey = iota + 1
	pathKey
)

func setHandler(ctx context.Context, handler Handler) context.Context {
	return context.WithValue(ctx, handlerKey, handler)
}

func getHandler(ctx context.Context) Handler {
	if handler, exist := ctx.Value(handlerKey).(Handler); exist {
		return handler
	}
	return nil
}

func setPath(ctx context.Context, path string) context.Context {
	return context.WithValue(ctx, pathKey, path)
}

func setEmptyPath(ctx context.Context) context.Context {
	return setPath(ctx, "")
}

func getPath(ctx context.Context, r *http.Request) string {
	if path, exist := ctx.Value(pathKey).(string); exist {
		return path
	}
	return r.URL.EscapedPath()
}
