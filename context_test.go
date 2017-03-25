package dispatch

import (
	"net/http/httptest"
	"testing"
)

func TestParam(t *testing.T) {
	key := "param1"
	expected := ""
	r := httptest.NewRequest("", "/", nil)
	actual := Param(r, key)
	if expected != actual {
		t.Fatalf("registered params: expected empty, actual %v", actual)
	}

	expected = "value1"
	r = setParam(r, key, expected)
	actual = Param(r, key)
	if expected != actual {
		t.Fatalf("registered param: expected %v, actual %v", expected, actual)
	}

	expected = "value2"
	r = setParam(r, key, expected)
	actual = Param(r, key)
	if expected != actual {
		t.Fatalf("registered params: expected %v, actual %v", expected, actual)
	}
}
