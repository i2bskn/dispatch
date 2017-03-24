package dispatch

import (
	"net/http"
)

// Route contains handler and routing information.
type Route struct {
	pattern string
	origin  http.Handler
	handler http.Handler
	method  HTTPMethod
	mux     *Mux
}

func newRoute(pattern string, origin http.Handler, mux *Mux) *Route {
	route := &Route{
		pattern: pattern,
		origin:  origin,
		method:  MethodAny,
		mux:     mux,
	}
	route.buildHandler()
	return route
}

// Methods is change the allow HTTP methods.
func (rt *Route) Methods(m HTTPMethod) *Route {
	rt.method = m
	return rt
}

func (rt *Route) isAcceptMethod(m string) bool {
	if (rt.method & parseMethod(m)) == 0 {
		return false
	}
	return true
}

func (rt *Route) buildHandler() {
	rt.handler = rt.origin
	for i := len(rt.mux.middleware) - 1; i >= 0; i-- {
		rt.handler = rt.mux.middleware[i](rt.handler)
	}
}

// HTTPMethod is type of HTTP method.
type HTTPMethod uint16

// These define HTTP methods that each Route to allow.
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

func parseMethod(method string) HTTPMethod {
	switch method {
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
