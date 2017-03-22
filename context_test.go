package tensile

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeRequest() *http.Request {
	return httptest.NewRequest("", "/", nil)
}

func TestParam(t *testing.T) {
	key := "key"
	expected := "value"
	r := setParam(fakeRequest(), key, expected)
	actual := Param(r, key)
	if expected != actual {
		t.Fatalf("registerd params: expected %v, actual %v", expected, actual)
	}
}

func TestSetParam(t *testing.T) {
	key := "key"
	expected := "value"
	r := setParam(fakeRequest(), key, expected)
	ctx := r.Context()
	p, ok := ctx.Value(paramKey).(map[string]string)
	if !ok {
		t.Fatal("params is not registered")
	}

	actual := p[key]
	if expected != actual {
		t.Fatalf("registerd params: expected %v, actual %v", expected, actual)
	}
}
