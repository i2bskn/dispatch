package pygmy

import (
	"context"
	"net/http"
)

type CompatibleHandler struct {
	Handler
}

func Compatible(handler Handler) http.Handler {
	return CompatibleHandler{
		Handler: handler,
	}
}

func (c CompatibleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler.ServeHTTP(nil, w, r)
}

type HTTPHandlerWrapper struct {
	http.Handler
}

func (h HTTPHandlerWrapper) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}
