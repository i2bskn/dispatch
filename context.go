package dispatch

import (
	"context"
	"net/http"
)

type contextKey int

const (
	paramKey contextKey = iota + 1
)

// Param returns path parameter from context.
func Param(r *http.Request, key string) string {
	ctx := r.Context()
	if p, ok := ctx.Value(paramKey).(map[string]string); ok {
		return p[key]
	}
	return ""
}

func setParam(r *http.Request, key, value string) *http.Request {
	ctx := r.Context()
	p, ok := ctx.Value(paramKey).(map[string]string)
	if !ok {
		p = make(map[string]string)
	}
	p[key] = value
	ctx = context.WithValue(ctx, paramKey, p)
	return r.WithContext(ctx)
}
