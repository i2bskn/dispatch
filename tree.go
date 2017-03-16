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
)

const paramEntrySize = 5

type node struct {
	label     string
	children  []*node
	kind      nodeKind
	entries   []*Entry
	paramName string
}

func (n *node) add(pattern string, e *Entry) bool {
	if n.kind == param {
		if len(pattern) < len(n.label) {
			return false
		}

		i := strings.IndexByte(pattern, '/')
		p := pattern
		if i >= 0 {
			p = pattern[:i]
		}
		if n.label != p {
			return false
		}

		if i == -1 {
			n.entries = append(n.entries, e)
			return true
		}

		pattern = pattern[i:]
	} else if n.kind != root {
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

	for _, child := range n.children {
		if child.add(pattern, e) {
			return true
		}
	}

	n.createChild(pattern, e)
	return true
}

func (n *node) createChild(pattern string, e *Entry) {
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
				child.createChild(pattern, e)
			}
		case ':':
			child.kind = param
			i := strings.IndexByte(pattern, '/')
			if i > 0 {
				child.label = pattern[:i]
				child.paramName = pattern[1:i]
				child.createChild(pattern[i:], e)
			} else {
				child.label = pattern
				child.paramName = pattern[1:]
				child.entries = append(child.entries, e)
			}
		default:
			child.kind = static
			i := strings.IndexByte(pattern, '/')
			if i > 0 {
				child.label = pattern[:i]
				child.createChild(pattern[i:], e)
			} else {
				child.label = pattern
				child.entries = append(child.entries, e)
			}
		}
		n.children = append(n.children, child)
	}
}

func (n *node) split(i int) {
	if 0 < i && i < len(n.label) {
		child := new(node)
		child.label = n.label[i:]
		child.kind = n.kind
		n.label = n.label[:i]
		n.children, child.children = child.children, n.children
		n.children = append(n.children, child)
		n.entries, child.entries = child.entries, n.entries
	}
}

func (n *node) match(path string, r *http.Request) *Entry {
	if n.kind == param {
		i := strings.IndexByte(path, '/')
		if i >= 0 {
			path = path[i:]
		} else {
			for _, e := range n.entries {
				if e.isAcceptMethod(r.Method) {
					return e
				}
			}
		}
	} else if n.kind != root {
		i := len(n.label)
		if len(path) < i || path[:i] != n.label {
			return nil
		}

		if n.kind == slash && len(n.entries) > 0 && len(r.URL.Path) > len(path) {
			for _, e := range n.entries {
				if e.isAcceptMethod(r.Method) {
					r.URL.Path = path
					return e
				}
			}
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

	children := make([]*node, 0, paramEntrySize)
	for _, child := range n.children {
		if child.kind == param {
			children = append(children, child)
			continue
		}

		if e := child.match(path, r); e != nil {
			return e
		}
	}

	for _, child := range children {
		if e := child.match(path, r); e != nil {
			return e
		}
	}

	return nil
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
