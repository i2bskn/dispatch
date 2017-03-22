package tensile

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

var tensileMethods = []HTTPMethod{
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

func fakeHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func TestIsAcceptMethod(t *testing.T) {
	entry := newEntry("/", fakeHandlerFunc())
	for _, method := range methods {
		if !entry.isAcceptMethod(method) {
			t.Fatalf("default should be accept all methods but %s not accept", method)
		}
	}
}

func TestParseMethod(t *testing.T) {
	for i := 0; i < len(methods); i++ {
		actual := parseMethod(methods[i])
		if actual != tensileMethods[i] {
			t.Fatalf("parseMethod(%s): expected %v, actual %v", methods[i], tensileMethods[i], actual)
		}
	}
}
