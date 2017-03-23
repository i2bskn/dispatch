package tensile

import (
	"net/http"
	"path"
	"sync"
)

// Mux is http.ServeMux compatible HTTP multiplexer.
type Mux struct {
	mu              sync.RWMutex
	entries         *node
	middleware      []func(http.Handler) http.Handler
	NotFoundHandler http.Handler
}

// New allocates and returns a new Mux.
func New() *Mux {
	return &Mux{
		NotFoundHandler: http.NotFoundHandler(),
	}
}

// Handle registers handler and returns Entry.
func (mux *Mux) Handle(pattern string, h http.Handler) *Entry {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}

	if h == nil {
		panic("http: nil handler")
	}

	if mux.entries == nil {
		mux.entries = new(node)
	}

	p := cleanPath(pattern)
	e := newEntry(p, h, mux)
	mux.entries.add(p, e)
	return e
}

// HandleFunc registers handler function and returns Entry.
func (mux *Mux) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) *Entry {
	return mux.Handle(pattern, http.HandlerFunc(h))
}

// ServeHTTP dispatches matching requests to handlers.
func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h, r := mux.Handler(r)
	h.ServeHTTP(w, r)
}

// Handler returns a handler to dispatch from request.
func (mux *Mux) Handler(r *http.Request) (http.Handler, *http.Request) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if e, r := mux.entries.match(r.URL.Path, r); e != nil {
		return e.handler, r
	}

	return mux.NotFoundHandler, r
}

// Use registers middleware.
func (mux *Mux) Use(middleware ...func(http.Handler) http.Handler) {
	mux.middleware = append(mux.middleware, middleware...)
	mux.entries.traverse(func(e *Entry) {
		e.buildHandler()
	})
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}

	if p[0] != '/' {
		p = "/" + p
	}

	np := path.Clean(p)
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

func buildHandler(h http.Handler, middleware ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}
	return h
}
