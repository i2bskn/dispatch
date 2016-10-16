package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/i2bskn/pygmy"
)

type LoggerFunc func(context.Context, *http.Request, time.Duration)

func NewLogger() pygmy.MiddlewareFunc {
	logger := log.New(os.Stdout, "", log.LstdFlags)
	f := func(ctx context.Context, r *http.Request, p time.Duration) {
		logger.Printf("%s %s %s\n", r.Method, r.URL.EscapedPath(), p.String())
	}
	return MakeLogger(f)
}

func MakeLogger(f LoggerFunc) pygmy.MiddlewareFunc {
	return func(h pygmy.Handler) pygmy.Handler {
		return pygmy.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			h.ServeHTTP(ctx, w, r)
			latency := time.Since(start)
			f(ctx, r, latency)
		})
	}
}
