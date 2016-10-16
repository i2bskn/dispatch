package pygmy

import (
	"context"
	"net/http"
	"path"
)

type Handler interface {
	ServeHTTP(context.Context, http.ResponseWriter, *http.Request)
}

type HandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (f HandlerFunc) ServeHTTP(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	f(ctx, w, r)
}

type MiddlewareFunc func(Handler) Handler

type Mux struct {
	root       bool
	ctx        context.Context
	routes     []*Route
	middleware []MiddlewareFunc
}

func New() *Mux {
	return &Mux{
		root: true,
		ctx:  context.Background(),
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
		obj := getShare(c)
		h := obj.handler
		for i := len(obj.middleware) - 1; i >= 0; i-- {
			h = obj.middleware[i](h)
		}
		h.ServeHTTP(c, w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (mux *Mux) Use(middleware ...MiddlewareFunc) {
	mux.middleware = append(mux.middleware, middleware...)
}

func (mux *Mux) Compatible() http.Handler {
	return Compatible(mux)
}

func (mux *Mux) belonging() {
	mux.root = false
}

func (mux *Mux) context(ctx context.Context) context.Context {
	if mux.root && mux.ctx != nil {
		return mux.ctx
	}

	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func (mux *Mux) match(ctx context.Context, r *http.Request) (context.Context, bool) {
	cc := ctx
	obj := getShare(cc)
	if obj == nil {
		path := r.URL.EscapedPath()
		obj = &share{path: cleanPath(path)}
		if obj.path != path {
			obj.handler = HTTPHandlerWrapper{
				http.RedirectHandler(obj.path, http.StatusMovedPermanently),
			}
			return setShare(cc, obj), true
		}
	}
	obj.middleware = append(obj.middleware, mux.middleware...)
	cc = setShare(cc, obj)

	for _, route := range mux.routes {
		mc, ok := route.match(cc, r)
		if ok {
			return mc, true
		}
	}

	return ctx, false
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

type share struct {
	path       string
	handler    Handler
	middleware []MiddlewareFunc
}

func (s *share) foundRoute() {
	s.path = ""
}
