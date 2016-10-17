package pygmy

import (
	"context"
	"net/http"
)

// CompatibleHandler is pygmy.Handler wrapper for the compatible to http.Handler
type CompatibleHandler struct {
	Handler
}

// Compatible is returns http.Handler by pygmy.Handler
func Compatible(handler Handler) http.Handler {
	return CompatibleHandler{
		Handler: handler,
	}
}

func (c CompatibleHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Handler.ServeHTTP(nil, w, r)
}

// HTTPHandlerWrapper is http.Handler wrapper for the compatible to pygmy.Handler
type HTTPHandlerWrapper struct {
	http.Handler
}

func (h HTTPHandlerWrapper) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	h.Handler.ServeHTTP(w, r)
}
