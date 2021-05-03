package routing

import (
	"net/http"
	"regexp"
)

const (
	catchAllExpression = "^.*$"
)

type node struct {
	prefix   string
	handler  http.HandlerFunc
	child    *node
	parent   *node
	sibling  *node
	t        int
	stops    map[byte]*node
	regexp   *regexp.Regexp
	w        int
	matchers []matcher
}

func (n *node) match(request *http.Request) bool {
	if n.handler == nil {
		return false
	}

	for _, m := range n.matchers {
		if ok, _ := m(request); !ok {
			return false
		}
	}
	return true
}

func (n *node) isCatchAll() bool {
	return n.regexpToString() == catchAllExpression
}

func (n *node) regexpEquals(o *node) bool {
	return n.regexpToString() == o.regexpToString()
}

func (n *node) regexpToString() string {
	if n.t != nodeTypeDynamic {
		return ""
	}

	if n.regexp == nil {
		return ""
	}
	return n.regexp.String()
}

func (n *node) hasParameters() bool {

	for _, m := range n.matchers {
		if _, hostLeaf := m(nil); hostLeaf.hasParameters() {
			return true
		}
	}

	parent := n
	for parent != nil {
		if parent.t == nodeTypeDynamic {
			return true
		}

		parent = parent.parent
	}

	return false
}
