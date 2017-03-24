package dispatch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type middlewareTest struct {
	wrapped bool
}

func fakeMiddleware(m *middlewareTest) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		m.wrapped = true
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	}
}

func fakeHandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
}

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
	mux.Handle("", fakeHandlerFunc())
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

func TestServeHTTP__notFound(t *testing.T) {
	mux := New()
	mux.Handle("/", fakeHandlerFunc())

	r := httptest.NewRequest(http.MethodGet, "/notfound", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, r)

	res := w.Result()
	if res.StatusCode != http.StatusNotFound {
		t.Fatalf("NotFound: expected status code %v, actual %v", http.StatusNotFound, res.StatusCode)
	}
}

func TestUse(t *testing.T) {
	mt := new(middlewareTest)
	mux := New()
	mux.Handle("/", fakeHandlerFunc())

	mux.Use(fakeMiddleware(mt))
	if !mt.wrapped {
		t.Fatal("middleware is not applied")
	}
}

func TestCleanPath(t *testing.T) {
	expected := "/"
	actual := cleanPath("")
	if expected != actual {
		t.Fatalf("cleanPath(\"\"): expected %s, actual %s", expected, actual)
	}

	expected = "/example"
	actual = cleanPath("example")
	if expected != actual {
		t.Fatalf("cleanPath(\"example\"): expected %s, actual %s", expected, actual)
	}

	expected = "/"
	actual = cleanPath("/../..")
	if expected != actual {
		t.Fatalf("cleanPath(\"/../..\"): expected %s, actual %s", expected, actual)
	}

	expected = "/example/"
	actual = cleanPath(expected)
	if expected != actual {
		t.Fatalf("cleanPath(expected): expected %s, actual %s", expected, actual)
	}
}
