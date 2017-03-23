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
	mux := New()
	entry := newEntry("/", fakeHandlerFunc(), mux)
	for _, method := range methods {
		if !entry.isAcceptMethod(method) {
			t.Fatalf("default should be accept all methods but %s not accept", method)
		}
	}

	entry.Methods(MethodPost | MethodPut)
	for _, method := range methods {
		if method == http.MethodPost || method == http.MethodPut {
			continue
		}

		if entry.isAcceptMethod(method) {
			t.Fatalf("%s match to (MethodPost | MethodPut)", method)
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

	expected := MethodUnknown
	actual := parseMethod("UNKNOWN")
	if expected != actual {
		t.Fatalf("parseMethod(\"UNKNOWN\"): expected %v, actual %v", expected, actual)
	}
}
