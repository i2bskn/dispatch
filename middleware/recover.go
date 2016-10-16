package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/i2bskn/pygmy"
)

type RecoverFunc func(context.Context, http.ResponseWriter, string)

func NewRecover() pygmy.MiddlewareFunc {
	logger := log.New(os.Stderr, "", log.LstdFlags)
	f := func(ctx context.Context, w http.ResponseWriter, err string) {
		logger.Printf("Recover: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, http.StatusText(http.StatusInternalServerError))
	}
	return MakeRecover(f)
}

func MakeRecover(f RecoverFunc) pygmy.MiddlewareFunc {
	return func(h pygmy.Handler) pygmy.Handler {
		return pygmy.HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err != nil {
					f(ctx, w, fmt.Sprint(err))
				}
			}()

			h.ServeHTTP(ctx, w, r)
		})
	}
}
