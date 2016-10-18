package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/i2bskn/pygmy"
)

func fakeRecoverFunc(o io.Writer) RecoverFunc {
	return func(ctx context.Context, w http.ResponseWriter, err string) {
		fmt.Fprint(o, err)
	}
}

func TestNewRecover(t *testing.T) {
	var _ pygmy.MiddlewareFunc = NewRecover()
}

func TestRecover(t *testing.T) {
	want := "PANIC"
	buf := new(bytes.Buffer)
	mux := pygmy.New()
	mux.Use(MakeRecover(fakeRecoverFunc(buf)))
	mux.HandleFunc("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		panic(want)
	})

	ts := httptest.NewServer(mux.Compatible())
	defer ts.Close()

	_, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("http.Get: %v", err)
	}

	if buf.String() != want {
		t.Fatalf("got: %#v\nwant: %#v", buf.String(), want)
	}
}
