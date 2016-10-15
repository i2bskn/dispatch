package pygmy

import (
	"context"
	"net/http"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

type Mux struct {
	root   bool
	ctx    context.Context
	routes []*Route
}

func New() *Mux {
	return &Mux{
		root:   true,
		ctx:    context.Background(),
		routes: make([]*Route, 0),
	}
}

func (mux *Mux) Handle(path string, h Handler) *Route {
	route := newRoute(path, h)
	mux.routes = append(mux.routes, route)
	return route
}

func (mux *Mux) HandleFunc(path string, h func(context.Context, http.ResponseWriter, *http.Request)) *Route {
	return mux.Handle(path, HandlerFunc(h))
}

func (mux *Mux) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	c, ok := mux.match(mux.context(ctx), r)
	if ok {
		handler := getHandler(c)
		handler.ServeHTTP(c, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (mux *Mux) Compatible() http.Handler {
	return Compatible(mux)
}

func (mux *Mux) context(ctx context.Context) context.Context {
	if mux.root && mux.ctx != nil {
		return mux.ctx
	}

	if ctx == nil {
		context.Background()
	}

	return ctx
}

func (mux *Mux) match(ctx context.Context, r *http.Request) (context.Context, bool) {
	for _, route := range mux.routes {
		c, ok := route.match(ctx, r)
		if ok {
			return c, true
		}
	}
	return ctx, false
}
