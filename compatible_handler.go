package pygmy

import (
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
