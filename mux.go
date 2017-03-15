package tensile

import (
	"net/http"
	"path"
	"sync"
)

type Mux struct {
	mu         sync.RWMutex
	m          *node
	middleware []MiddlewareFunc
}

func New() *Mux {
	return new(Mux)
}

func (mux *Mux) Handle(pattern string, h http.Handler) *Entry {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if pattern == "" {
		panic("http: invalid pattern " + pattern)
	}

	if h == nil {
		panic("http: nil handler")
	}

	if mux.m == nil {
		mux.m = new(node)
	}

	p := cleanPath(pattern)
	e := newEntry(p, h)
	mux.m.add(p, e)
	return e
}

func (mux *Mux) HandleFunc(pattern string, h func(http.ResponseWriter, *http.Request)) *Entry {
	return mux.Handle(pattern, http.HandlerFunc(h))
}

func (mux *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h := mux.Handler(r)
	for i := len(mux.middleware) - 1; i >= 0; i-- {
		h = mux.middleware[i](h)
	}
	h.ServeHTTP(w, r)
}

func (mux *Mux) Handler(r *http.Request) http.Handler {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	if e := mux.m.match(r.URL.Path, r); e != nil {
		return e.h
	}

	return http.NotFoundHandler()
}

func (mux *Mux) Use(middleware ...MiddlewareFunc) {
	mux.middleware = append(mux.middleware, middleware...)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}

	if p[0] != '/' {
		p = "/" + p
	}

	return path.Clean(p)
}
