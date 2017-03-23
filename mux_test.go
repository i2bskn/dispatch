package tensile

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleFunc(t *testing.T) {
	mux := New()
	if mux.entries != nil {
		t.Fatal("Mux just after initialization has no entries")
	}

	expected := "body"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	})
	if mux.entries == nil {
		t.Fatal("Mux always has entries after handler registration")
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	res := w.Result()
	body, _ := ioutil.ReadAll(res.Body)
	if string(body) != expected {
		t.Fatal("registered handler has not called")
	}
}

func TestHandle__invalidPattern(t *testing.T) {
	defer func() {
		recover()
	}()

	mux := New()
	mux.Handle("", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	t.Fatal("panic has not occurred with invalid pattern")
}

func TestHandle__emptyHandler(t *testing.T) {
	defer func() {
		recover()
	}()

	mux := New()
	mux.Handle("/", nil)
	t.Fatal("panic has not occurred with empty handler")
}
