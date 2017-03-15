package tensile

import (
	"net/http"
	"strings"
)

type nodeKind uint8

const (
	root nodeKind = iota
	slash
	static
	param
	any
)

type node struct {
	label    string
	childlen []*node
	kind     nodeKind
	entries  []*Entry
}

func (n *node) add(pattern string, e *Entry) bool {
	if n.kind != root {
		i := 0
		for i < min(len(n.label), len(pattern)) && n.label[i] == pattern[i] {
			i++
		}

		if i == 0 {
			return false
		}
		if len(n.label) > i && n.kind == static {
			n.split(i)
		}
		pattern = pattern[i:]
	}

	for _, child := range n.childlen {
		if child.add(pattern, e) {
			return true
		}
	}

	n.insertChild(pattern, e)
	return true
}

func (n *node) match(path string, r *http.Request) *Entry {
	if n.kind != root {
		i := len(n.label)
		if len(path) < i || path[:i] != n.label {
			return nil
		}

		if len(path) == i {
			for _, e := range n.entries {
				if e.isAcceptMethod(r.Method) {
					return e
				}
			}
		}

		path = path[i:]
	}

	for _, child := range n.childlen {
		if e := child.match(path, r); e != nil {
			return e
		}
	}

	return nil
}

func (n *node) insertChild(pattern string, e *Entry) {
	if len(pattern) > 0 {
		child := new(node)
		switch pattern[0] {
		case '/':
			child.label = pattern[:1]
			child.kind = slash
			pattern = pattern[1:]
			if len(pattern) == 0 {
				child.entries = append(child.entries, e)
			} else {
				child.insertChild(pattern, e)
			}
		case ':':
		default:
			child.kind = static
			i := strings.IndexByte(pattern, '/')
			if i > 0 {
				child.label = pattern[:i]
				child.insertChild(pattern[i:], e)
			} else {
				child.label = pattern
				child.entries = append(child.entries, e)
			}
		}
		n.childlen = append(n.childlen, child)
	}
}

func (n *node) split(i int) {
	if 0 < i && i < len(n.label) {
		child := new(node)
		child.label = n.label[i:]
		child.kind = n.kind
		n.label = n.label[:i]
		n.childlen, child.childlen = child.childlen, n.childlen
		n.childlen = append(n.childlen, child)
		n.entries, child.entries = child.entries, n.entries
	}
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
