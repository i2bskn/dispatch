package pygmy

import (
	"context"
	"net/http"
	"strings"
)

const (
	tokenSep    = '/'
	namedPrefix = ':'
	namedExp    = "*"
)

type matcher struct {
	raw        string
	tokens     []string
	names      map[int]string
	broadMatch bool
	static     bool
}

func newMatcher(path string, broadMatch bool) *matcher {
	tokens, names, static := buildTokens(path)
	return &matcher{
		raw:        path,
		tokens:     tokens,
		names:      names,
		broadMatch: broadMatch,
		static:     static,
	}
}

func (m *matcher) match(ctx context.Context, r *http.Request) (context.Context, bool) {
	c := ctx
	obj := getShare(c)
	rpath := obj.path
	if m.static {
		if m.broadMatch {
			if !strings.HasPrefix(rpath, m.raw) {
				return ctx, false
			}

			obj.path = rpath[len(m.raw):]
			return setShare(c, obj), true
		}

		if m.raw == rpath {
			obj.foundRoute()
			return setShare(c, obj), true
		} else {
			return ctx, false
		}
	}

	rtokenSize := strings.Count(rpath, string(tokenSep))
	tokenSize := len(m.tokens)
	if rtokenSize < tokenSize {
		return ctx, false
	}

	params := make(map[string]string)

	eop := len(rpath) - 1
	if rtokenSize > tokenSize {
		if m.broadMatch {
			s := 0
			for i := 1; i <= eop; i++ {
				if rpath[i] == tokenSep {
					s++
					if s == tokenSize {
						eop = i - 1
						obj.path = rpath[i:]
						c = setShare(c, obj)
						break
					}
				}
			}
		} else {
			return ctx, false
		}
	}

	fot := 0
	idx := 0
	for i := 1; i <= eop; i++ {
		if rpath[i] == tokenSep {
			rtoken := rpath[fot:i]
			if m.tokens[idx] == namedExp {
				params[m.names[idx]] = rtoken[1:]
			} else {
				if m.tokens[idx] != rtoken {
					return ctx, false
				}
			}
			fot = i
			idx++
		} else if i == eop {
			rtoken := rpath[fot:i]
			if m.tokens[idx] == namedExp {
				params[m.names[idx]] = rtoken[1:]
			} else {
				if m.tokens[idx] != rtoken {
					return ctx, false
				}
			}
		}
	}

	if !m.broadMatch {
		obj.foundRoute()
		c = setShare(c, obj)
	}
	return c, true
}

func buildTokens(path string) ([]string, map[int]string, bool) {
	tokens := make([]string, strings.Count(path, string(tokenSep)))
	names := make(map[int]string)
	static := true

	eop := len(path) - 1
	fot := 0
	idx := 0
	for i := 1; i <= eop; i++ {
		if path[i] == tokenSep {
			token := path[fot:i]
			if token[1] == namedPrefix {
				names[idx] = token[2:]
				token = namedExp
			}
			tokens[idx] = token
			fot = i
			idx++
		} else if i == eop {
			token := path[fot:]
			if token[1] == namedPrefix {
				names[idx] = token[2:]
				token = namedExp
			}
			tokens[idx] = token
		}
	}

	if len(names) > 0 {
		static = false
	}
	return tokens, names, static
}
