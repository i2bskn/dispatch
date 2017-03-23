package tensile

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func fakeRequest() *http.Request {
	return httptest.NewRequest("", "/", nil)
}

func TestParam__match(t *testing.T) {
	key := "key"
	expected := "value"
	r := setParam(fakeRequest(), key, expected)
	actual := Param(r, key)
	if expected != actual {
		t.Fatalf("registerd params: expected %v, actual %v", expected, actual)
	}
}

func TestParam__unmatch(t *testing.T) {
	key := "key"
	expected := ""
	r := fakeRequest()
	actual := Param(r, key)
	if expected != actual {
		t.Fatalf("registerd params: expected empty, actual %v", actual)
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

	expected = "update"
	r = setParam(r, key, expected)
	ctx = r.Context()
	p, ok = ctx.Value(paramKey).(map[string]string)
	if !ok {
		t.Fatal("params is not registered")
	}

	actual = p[key]
	if expected != actual {
		t.Fatalf("registerd params: expected %v, actual %v", expected, actual)
	}
}
