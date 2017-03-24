package dispatch

import (
	"net/http/httptest"
	"testing"
)

type routeTest struct {
	pattern          string
	requestPath      string
	pathAfterMatched string
	paramName        string
	paramValue       string
}

func fakeEntry(pattern string) *Entry {
	return &Entry{
		pattern: pattern,
		method:  MethodAny,
	}
}

func entryCount(tree *node) int {
	i := 0
	tree.traverse(func(e *Entry) {
		i++
	})
	return i
}

func testEntryCount(t *testing.T, tree *node, expected int) {
	actual := entryCount(tree)
	if expected != actual {
		t.Fatalf("number of entries unexpected: expected %d, actual %d", expected, actual)
	}
}

func testRouteMatch(t *testing.T, tree *node, tests []routeTest) {
	for _, tt := range tests {
		r := httptest.NewRequest("", tt.requestPath, nil)
		e, r := tree.match(r.URL.Path, r)
		if e == nil {
			t.Fatalf("route not found: %s", tt.requestPath)
		}

		if e.pattern != tt.pattern {
			t.Fatalf("pattern unmatch: expected %s, actual %s", tt.pattern, e.pattern)
		}

		if r.URL.Path != tt.pathAfterMatched {
			t.Fatalf("unexpected path after matched: expected %s, actual %s", tt.pathAfterMatched, r.URL.Path)
		}

		if len(tt.paramName) > 0 {
			actual := Param(r, tt.paramName)
			if tt.paramValue != actual {
				t.Fatalf("unexpected param value: path %s, name %s, value %s", tt.requestPath, tt.paramName, actual)
			}
		}
	}
}

func testNotFound(t *testing.T, tree *node) {
	r := httptest.NewRequest("", "/notfound", nil)
	e, r := tree.match(r.URL.Path, r)
	if e != nil {
		t.Fatalf("unexpected match route: %s", e.pattern)
	}
}

func TestTree(t *testing.T) {
	tree := new(node)
	tests := []routeTest{
		{"/abc", "/abc", "/abc", "", ""},
		{"/abc/", "/abc/def", "/def", "", ""},
		{"/aaa/:id", "/aaa/123", "/aaa/123", "id", "123"},
		{"/aaa/:id/bbb", "/aaa/456/bbb", "/aaa/456/bbb", "id", "456"},
		{"/aaa/bbb/ccc", "/aaa/bbb/ccc", "/aaa/bbb/ccc", "", ""},
		{"/", "/", "/", "", ""},
	}

	for _, tt := range tests {
		tree.add(tt.pattern, fakeEntry(tt.pattern))
	}

	testEntryCount(t, tree, len(tests))
	testRouteMatch(t, tree, tests)
	testNotFound(t, tree)
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
