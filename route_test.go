package dispatch

import (
	"net/http"
	"testing"
)

var methods = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

var pkgMethods = []HTTPMethod{
	MethodGet,
	MethodHead,
	MethodPost,
	MethodPut,
	MethodPatch,
	MethodDelete,
	MethodConnect,
	MethodOptions,
	MethodTrace,
}

func TestRoute(t *testing.T) {
	mt := new(middlewareTest)
	mux := New()
	mux.middleware = append(mux.middleware, fakeMiddleware(mt))
	route := newRoute("/", fakeHandlerFunc(), mux)

	if !mt.wrapped {
		t.Fatal("middleware is not applied")
	}

	for _, method := range methods {
		if !route.isAcceptMethod(method) {
			t.Fatalf("default should be accept all methods but %s not accept", method)
		}
	}

	expected := MethodPost | MethodPut
	route.Methods(expected)
	if route.method != expected {
		t.Fatalf("accept methods is not update: expected %v, actual %v", expected, route.method)
	}

	for _, method := range methods {
		if method == http.MethodPost || method == http.MethodPut {
			continue
		}

		if route.isAcceptMethod(method) {
			t.Fatalf("%s match to (MethodPost | MethodPut)", method)
		}
	}
}

func TestParseMethod(t *testing.T) {
	for i := 0; i < len(methods); i++ {
		actual := parseMethod(methods[i])
		if actual != pkgMethods[i] {
			t.Fatalf("parseMethod(%s): expected %v, actual %v", methods[i], pkgMethods[i], actual)
		}
	}

	expected := MethodUnknown
	actual := parseMethod("UNKNOWN")
	if expected != actual {
		t.Fatalf("parseMethod(\"UNKNOWN\"): expected %v, actual %v", expected, actual)
	}
}
