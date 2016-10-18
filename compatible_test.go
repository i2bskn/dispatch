package pygmy

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeCompatibleHandler(want string) http.Handler {
	mux := New()
	mux.HandleFunc("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, want)
	})
	return Compatible(mux)
}

func TestCompatible(t *testing.T) {
	var _ http.Handler = Compatible(New())
}

func TestCompatibleHandler(t *testing.T) {
	want := "Hello Pygmy!"
	ts := httptest.NewServer(fakeCompatibleHandler(want))
	defer ts.Close()

	r, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("ioutil.ReadAll: %v", err)
	}

	if want != string(data) {
		t.Fatalf("got: %#v\nwant: %#v", string(data), want)
	}
}

func TestHTTPHandlerWrapper(t *testing.T) {
	var _ Handler = HTTPHandlerWrapper{http.RedirectHandler("/", http.StatusMovedPermanently)}
}
