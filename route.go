package pygmy

import (
	"context"
	"net/http"
)

// Route contains handler and routing information.
type Route struct {
	leaf    bool
	handler Handler
	mux     *Mux
	method  HTTPMethod
	pattern *matcher
}

func newRoute(path string, handler Handler) *Route {
	route := &Route{
		leaf:   true,
		method: MethodAny,
	}

	switch h := handler.(type) {
	case *Mux:
		h.belonging()
		route.leaf = false
		route.mux = h
	default:
		route.handler = h
	}

	route.pattern = newMatcher(path, route.leaf)

	return route
}

// Verbs is change the allow HTTP methods.
func (rt *Route) Verbs(m HTTPMethod) *Route {
	rt.method = m
	return rt
}

func (rt *Route) match(ctx context.Context, r *http.Request) (context.Context, bool) {
	if !rt.isAcceptMethod(r.Method) {
		return ctx, false
	}

	c, ok := rt.pattern.match(ctx, r)
	if !ok {
		return ctx, false
	}

	if rt.leaf {
		obj := getShare(c)
		obj.handler = rt.handler
		c = setShare(c, obj)
	} else {
		c, ok = rt.mux.match(c, r)
		if !ok {
			return ctx, false
		}
	}
	return c, true
}

func (rt *Route) isAcceptMethod(m string) bool {
	if (rt.method & parseMethod(m)) == 0 {
		return false
	}
	return true
}

// HTTPMethod is type of HTTP method flag.
type HTTPMethod uint16

// These flags define HTTP methods that each Route to allow.
const (
	MethodGet HTTPMethod = 1 << iota
	MethodHead
	MethodPost
	MethodPut
	MethodPatch
	MethodDelete
	MethodConnect
	MethodOptions
	MethodTrace
	MethodAny = MethodGet | MethodHead | MethodPost | MethodPut | MethodPatch | MethodDelete |
		MethodConnect | MethodOptions | MethodTrace
	MethodUnknown = HTTPMethod(0)
)

func parseMethod(m string) HTTPMethod {
	switch m {
	case http.MethodGet:
		return MethodGet
	case http.MethodHead:
		return MethodHead
	case http.MethodPost:
		return MethodPost
	case http.MethodPut:
		return MethodPut
	case http.MethodPatch:
		return MethodPatch
	case http.MethodDelete:
		return MethodDelete
	case http.MethodConnect:
		return MethodConnect
	case http.MethodOptions:
		return MethodOptions
	case http.MethodTrace:
		return MethodTrace
	default:
		return MethodUnknown
	}
}
