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
	parent  *node
	sibling *node
	t       int
	stops   map[byte]*node
	regexp  *regexp.Regexp
	w       int
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
	parent := n
	for parent != nil {
		if parent.t == nodeTypeDynamic {
			return true
		}
		parent = parent.parent
	}

	return false
}

type nodeInterface interface {
	hasParameters() bool
	sibling() nodeInterface
	child() nodeInterface
	find(p string) nodeInterface
	combine(n nodeInterface) nodeInterface
	match(p string) bool
	handler() http.HandlerFunc
}

type nodeStatic struct {
	prefix      string
	handlerFunc http.HandlerFunc
	childNode   nodeInterface
	parentNode  nodeInterface
	siblingNode nodeInterface
	weight      int
}

func (ns *nodeStatic) hasParameters() bool {
	return false
}

func (ns *nodeStatic) sibling() nodeInterface {
	return ns.siblingNode
}

func (ns *nodeStatic) child() nodeInterface {
	return ns.childNode
}

func (ns *nodeStatic) handler() http.HandlerFunc {
	return ns.handlerFunc
}

func (ns *nodeStatic) find(p string) nodeInterface {
	if len(p) == 0 {
		return nil
	}

	pos := common(p, ns.prefix)
	if pos == 0 {
		return ns.sibling().find(p)
	}

	if pos == len(p) && len(p) == len(ns.prefix) {
		if nil != ns.handlerFunc {
			return ns
		}

		return nil
	}

	n := ns.child().find(p[pos:])
	if nil != n && nil != n.handler() {
		return n
	}

	for next := ns.sibling(); nil != next; next = next.sibling() {
		if n := next.find(p); nil != n {
			return n
		}
	}

	return nil
}

func (ns *nodeStatic) combine(n nodeInterface) nodeInterface {
	return nil
}

func (ns *nodeStatic) match(p string) bool {
	return false
}

type nodeDynamic struct {
	prefix        string
	handlerFunc   http.HandlerFunc
	parentNode    nodeInterface
	siblingNode   nodeInterface
	childrenNodes map[byte]nodeInterface
	regexp        *regexp.Regexp
	weight        int
}

func (nd *nodeDynamic) hasParameters() bool {
	return false
}

func (nd *nodeDynamic) sibling() nodeInterface {
	return nd.siblingNode
}

func (nd *nodeDynamic) child() nodeInterface {
	return nil
}

func (nd *nodeDynamic) handler() http.HandlerFunc {
	return nd.handlerFunc
}

func (nd *nodeDynamic) find(p string) nodeInterface {
	if len(p) == 0 {
		return nil
	}

	traversed := false
	for i := 0; i < len(p); i++ {

		if next, ok := nd.childrenNodes[p[i]]; ok {
			validExpression := true
			if nd.regexp != nil {
				validExpression = nd.regexp.MatchString(p[0:i])
			}
			if validExpression {
				traversed = true
				found := next.find(p[i:])
				if nil != found && nil != found.handler() {
					return found
				}
			}
		}

		if p[i] == '/' && !nd.isCatchAll() {
			return nd.sibling().find(p)
		}
	}

	if nil != nd.handler() && !traversed {
		validExpression := true
		if nd.regexp != nil {
			validExpression = nd.regexp.MatchString(p)
		}
		if validExpression {
			return nd
		}
	}

	return nd.sibling().find(p)
}

func (nd *nodeDynamic) combine(n nodeInterface) nodeInterface {
	return nil
}

func (nd *nodeDynamic) match(p string) bool {
	return false
}

func (nd *nodeDynamic) isCatchAll() bool {
	if nd.regexp == nil {
		return false
	}
	return nd.regexp.String() == catchAllExpression
}
