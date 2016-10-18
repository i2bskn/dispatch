package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/i2bskn/pygmy"
)

func fakeLoggerFunc(o io.Writer, want string) LoggerFunc {
	return func(ctx context.Context, r *http.Request, p time.Duration) {
		fmt.Fprint(o, want)
	}
}

func TestNewLogger(t *testing.T) {
	var _ pygmy.MiddlewareFunc = NewLogger()
}

func TestLogger(t *testing.T) {
	want := "example log"
	buf := new(bytes.Buffer)
	mux := pygmy.New()
	mux.Use(MakeLogger(fakeLoggerFunc(buf, want)))
	mux.HandleFunc("/", func(ctx context.Context, w http.ResponseWriter, r *http.Request) {})

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
