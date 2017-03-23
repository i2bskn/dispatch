package tensile

import (
	"net/http/httptest"
	"testing"
)

func TestParam(t *testing.T) {
	key := "param1"
	expected := "value1"
	r := setParam(httptest.NewRequest("", "/", nil), key, expected)

	actual := Param(r, key)
	if expected != actual {
		t.Fatalf("registerd param: expected %v, actual %v", expected, actual)
	}

	expected = ""
	actual = Param(r, "param2")
	if expected != actual {
		t.Fatalf("registerd params: expected empty, actual %v", actual)
	}

	expected = "value2"
	r = setParam(r, key, expected)
	actual = Param(r, key)
	if expected != actual {
		t.Fatalf("registerd params: expected %v, actual %v", expected, actual)
	}
}
