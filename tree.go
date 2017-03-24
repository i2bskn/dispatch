package dispatch

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

const paramChildrenSize = 10

type node struct {
	label     string
	children  []*node
	kind      nodeKind
	routes    []*Route
	paramName string
}

func (n *node) add(pattern string, route *Route) bool {
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
			n.routes = append(n.routes, route)
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

		if len(pattern) == 0 {
			n.routes = append(n.routes, route)
			return true
		}
	}

	for _, child := range n.children {
		if child.add(pattern, route) {
			return true
		}
	}

	n.createChild(pattern, route)
	return true
}

func (n *node) createChild(pattern string, route *Route) {
	if len(pattern) > 0 {
		child := new(node)
		switch pattern[0] {
		case '/':
			child.label = pattern[:1]
			child.kind = slash
			pattern = pattern[1:]
			if len(pattern) == 0 {
				child.routes = append(child.routes, route)
			} else {
				child.createChild(pattern, route)
			}
		case ':':
			child.kind = param
			i := strings.IndexByte(pattern, '/')
			if i > 0 {
				child.label = pattern[:i]
				child.paramName = pattern[1:i]
				child.createChild(pattern[i:], route)
			} else {
				child.label = pattern
				child.paramName = pattern[1:]
				child.routes = append(child.routes, route)
			}
		default:
			child.kind = static
			i := strings.IndexByte(pattern, '/')
			if i > 0 {
				child.label = pattern[:i]
				child.createChild(pattern[i:], route)
			} else {
				child.label = pattern
				child.routes = append(child.routes, route)
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
		n.routes, child.routes = child.routes, n.routes
	}
}

func (n *node) match(path string, r *http.Request) (*Route, *http.Request) {
	if n.kind == param {
		i := strings.IndexByte(path, '/')
		if i >= 0 {
			r = setParam(r, n.paramName, path[:i])
			path = path[i:]
		} else {
			for _, route := range n.routes {
				if route.isAcceptMethod(r.Method) {
					r = setParam(r, n.paramName, path)
					return route, r
				}
			}
		}
	} else if n.kind != root {
		if n.label == path {
			for _, route := range n.routes {
				if route.isAcceptMethod(r.Method) {
					return route, r
				}
			}
		}

		i := len(n.label)
		if len(path) < i || path[:i] != n.label {
			return nil, r
		}

		if n.kind == slash && len(n.routes) > 0 && len(r.URL.Path) > len(path) {
			for _, route := range n.routes {
				if route.isAcceptMethod(r.Method) {
					r.URL.Path = path
					return route, r
				}
			}
		}

		path = path[i:]
	}

	children := make([]*node, 0, paramChildrenSize)
	for _, child := range n.children {
		if child.kind == param {
			children = append(children, child)
			continue
		}

		if route, r := child.match(path, r); route != nil {
			return route, r
		}
	}

	for _, child := range children {
		if route, r := child.match(path, r); route != nil {
			return route, r
		}
	}

	return nil, r
}

func (n *node) traverse(f func(*Route)) {
	for _, route := range n.routes {
		f(route)
	}

	for _, child := range n.children {
		child.traverse(f)
	}
}

func min(a, b int) int {
	if a >= b {
		return b
	}
	return a
}
