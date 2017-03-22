package tensile

import (
	"net/http"
	"testing"
)

func fakeHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

func TestIsAcceptMethod(t *testing.T) {
	entry := newEntry("/", fakeHandlerFunc())
	methods := []string{
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

	for _, method := range methods {
		if !entry.isAcceptMethod(method) {
			t.Fatalf("default should be accept all methods but %s not accept", method)
		}
	}
}
