package routing

import (
	"net/http"
	"regexp"
)

const (
	catchAllExpression = "^.*$"
)

type node struct {
	prefix  string
	handler http.HandlerFunc
	child   *node
	sibling *node
	t       int
	stops   map[byte]*node
	regexp  *regexp.Regexp
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