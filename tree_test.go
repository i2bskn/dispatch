package tensile

import (
	"net/http/httptest"
	"testing"
)

func TestTree(t *testing.T) {
	mux := New()
	patterns := []string{"/", "/abc", "/abc/", "/aaa/:name", "/aaa/bbb/ccc"}
	for _, pattern := range patterns {
		mux.Handle(pattern, fakeHandlerFunc())
	}

	entryCount := 0
	mux.entries.traverse(func(e *Entry) {
		entryCount++
	})
	if entryCount != len(patterns) {
		t.Fatalf("entryCount: expected %d, actual %d", len(patterns), entryCount)
	}

	for _, pattern := range patterns {
		r := httptest.NewRequest("", pattern, nil)
		e, r := mux.entries.match(pattern, r)
		if e == nil {
			t.Fatalf("node.match(%s): not found", pattern)
		}

		if e.pattern != pattern {
			t.Fatalf("node.match(%s): pattern unmatch: entry.pattern %s", pattern, e.pattern)
		}
	}

	r := httptest.NewRequest("", "/notfound", nil)
	e, r := mux.entries.match(r.URL.Path, r)
	if e != nil {
		t.Fatalf("node.match(\"/notfound\"): actual %s", e.pattern)
	}

	r = httptest.NewRequest("", "/aaa/123", nil)
	e, r = mux.entries.match(r.URL.Path, r)
	if e == nil {
		t.Fatal("param routing did not match")
	}
	if e.pattern != "/aaa/:name" {
		t.Fatalf("node.match(\"/aaa/123\"): actual %s", e.pattern)
	}
	actual := Param(r, "name")
	if actual != "123" {
		t.Fatalf("Param(r, \"name\"): expected 123, actual %s", actual)
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		in       [2]int
		expected int
	}{
		{[2]int{1, 2}, 1},
		{[2]int{2, 2}, 2},
		{[2]int{3, 2}, 2},
	}

	for _, tt := range tests {
		actual := min(tt.in[0], tt.in[1])
		if actual != tt.expected {
			t.Fatalf("min(%d, %d): expected %d, actual %d", tt.in[0], tt.in[1], tt.expected, actual)
		}
	}
}
