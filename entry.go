package dispatch

import (
	"net/http"
)

// Entry contains handler and routing information.
type Entry struct {
	pattern string
	origin  http.Handler
	handler http.Handler
	method  HTTPMethod
	mux     *Mux
}

func newEntry(pattern string, origin http.Handler, mux *Mux) *Entry {
	e := &Entry{
		pattern: pattern,
		origin:  origin,
		method:  MethodAny,
		mux:     mux,
	}
	e.buildHandler()
	return e
}

// Methods is change the allow HTTP methods.
func (e *Entry) Methods(m HTTPMethod) *Entry {
	e.method = m
	return e
}

func (e *Entry) isAcceptMethod(m string) bool {
	if (e.method & parseMethod(m)) == 0 {
		return false
	}
	return true
}

func (e *Entry) buildHandler() {
	e.handler = e.origin
	for i := len(e.mux.middleware) - 1; i >= 0; i-- {
		e.handler = e.mux.middleware[i](e.handler)
	}
}

// HTTPMethod is type of HTTP method flag.
type HTTPMethod uint16

// These flags define HTTP methods that each Entry to allow.
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
