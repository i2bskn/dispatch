package pygmy

import (
	"context"
	"net/http"
	"testing"
)

func fakeHandlerFunc() HandlerFunc {
	return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {})
}

func TestNewRoute__handler(t *testing.T) {
	mux := New()
	route := newRoute("/", mux)
	if route.mux == nil {
		t.Fatal("initialized Route with mux should have mux attribute")
	}

	if route.handler != nil {
		t.Fatal("initialized Route with mux should not have handler attribute")
	}

	if route.leaf {
		t.Fatal("initialized Route with mux should not be leaf node")
	}

	if route.method != MethodAny {
		t.Fatal("got: %#v\nwant: %#v", route.method, MethodAny)
	}
}

func TestNewRoute__handlerFunc(t *testing.T) {
	route := newRoute("/", fakeHandlerFunc())
	if route.mux != nil {
		t.Fatal("initialized Route without mux should not have mux attribute")
	}

	if route.handler == nil {
		t.Fatal("initialized Route without mux should have handler attribute")
	}

	if !route.leaf {
		t.Fatal("initialized Route without mux should be leaf node")
	}

	if route.method != MethodAny {
		t.Fatalf("got: %#v\nwant: %#v", route.method, MethodAny)
	}
}

func TestIsAcceptMethod(t *testing.T) {
	route := newRoute("/", fakeHandlerFunc())
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
		if !route.isAcceptMethod(method) {
			t.Fatalf("default should be accept all methods but %s not accept", method)
		}
	}

	route.Verbs(MethodPost | MethodPut)
	for _, method := range methods {
		if method == http.MethodPost || method == http.MethodPut {
			continue
		}

		if route.isAcceptMethod(method) {
			t.Fatalf("should not be accept other than the accepted methods", method)
		}
	}
}
