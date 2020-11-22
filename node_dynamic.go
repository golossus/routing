package routing

import (
	"net/http"
	"regexp"
)

const (
	catchAllExpression = "^.*$"
	catchAllSeparator  = '/'
)

type nodeDynamic struct {
	prefix        string
	handlerFunc   http.HandlerFunc
	parentNode    nodeInterface
	siblingNode   nodeInterface
	childrenNodes map[byte]nodeInterface
	regexp        *regexp.Regexp
	weight        int
}

func newNodeDynamic(prefix string, rex *regexp.Regexp, h http.HandlerFunc) *nodeDynamic {
	return &nodeDynamic{
		prefix:        prefix,
		handlerFunc:   h,
		parentNode:    nil,
		siblingNode:   nil,
		childrenNodes: make(map[byte]nodeInterface),
		regexp:        rex,
		weight:        0,
	}
}

func (nd *nodeDynamic) hasParameters() bool {
	return true
}

func (nd *nodeDynamic) setParent(parent nodeInterface) {
	nd.parentNode = parent
}

func (nd *nodeDynamic) setSibling(sibling nodeInterface) {
	nd.siblingNode = sibling
}

func (nd *nodeDynamic) getWeight() int {
	return nd.weight
}

func (nd *nodeDynamic) getPrefix() string {
	return nd.prefix
}

func (nd *nodeDynamic) getParent() nodeInterface {
	return nd.parentNode
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

func (nd *nodeDynamic) setHandler(handlerFunc http.HandlerFunc) {
	nd.handlerFunc = handlerFunc
}

func (nd *nodeDynamic) find(p string) nodeInterface {
	if len(p) == 0 {
		return nil
	}

	traversed := false
	for i := 0; i < len(p); i++ {

		if next, ok := nd.childrenNodes[p[i]]; ok {
			if nd.matches(p[0:i]) {
				traversed = true
				found := next.find(p[i:])
				if nil != found {
					return found
				}
			}
		}

		if p[i] == catchAllSeparator && !nd.catchesAll() {
			if nd.siblingNode == nil {
				return nil
			}
			return nd.siblingNode.find(p)
		}
	}

	if nil != nd.handlerFunc && !traversed && nd.matches(p) {
		return nd
	}

	if nd.siblingNode == nil {
		return nil
	}

	return nd.siblingNode.find(p)
}

func (nd *nodeDynamic) merge(ni nodeInterface) nodeInterface {
	if ni == nil {
		return nd
	}

	n, ok := ni.(*nodeDynamic)
	if !ok {
		return ni.merge(nd)
	}

	if nd.equals(n) {
		for _, child := range n.childrenNodes {
			nd.addChild(child)
		}

		if n.handlerFunc != nil {
			nd.handlerFunc = n.handlerFunc
		}

		return nd
	}

	return nd.addSibling(n)
}

func (nd *nodeDynamic) addSibling(s nodeInterface) nodeInterface {
	if s == nil {
		return nd
	}

	if nd.siblingNode == nil {
		s.setParent(nd.parentNode)
		nd.siblingNode = s
		return nd
	}

	nd.siblingNode = nd.siblingNode.merge(s)

	return nd
}

func (nd *nodeDynamic) addChild(child nodeInterface) nodeInterface {
	if child == nil {
		return nd
	}

	k := child.getPrefix()[0]
	if nk, ok := nd.childrenNodes[k]; ok {
		nd.childrenNodes[k] = nk.merge(child)
		return nd
	}

	child.setParent(nd)
	nd.childrenNodes[k] = child

	return nd
}

func (nd *nodeDynamic) equals(o *nodeDynamic) bool {
	return nd.prefix == o.prefix && nd.regexpToString() == o.regexpToString()
}

func (nd *nodeDynamic) regexpToString() string {
	if nd.regexp == nil {
		return ""
	}
	return nd.regexp.String()
}

// matches returns bool if node matches a given string
func (nd *nodeDynamic) matches(p string) bool {
	if nd.regexp == nil {
		return true
	}

	return nd.regexp.MatchString(p)
}

// catchesAll returns whether the node catches all bytes of a prefix string
func (nd *nodeDynamic) catchesAll() bool {
	if nd.regexp == nil {
		return false
	}

	return nd.regexp.String() == catchAllExpression
}
