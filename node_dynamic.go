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

func (nd *nodeDynamic) combine(ni nodeInterface) nodeInterface {
	if ni == nil {
		return nd
	}

	n, ok := ni.(*nodeDynamic)
	if !ok {
		return ni.combine(nd)
	}

	if nd.equals(n) {
		for k, next1 := range nd.childrenNodes {
			next2, _ := n.childrenNodes[k]
			n.childrenNodes[k] = next1.combine(next2)
		}

		nd.childrenNodes = n.childrenNodes
		if n.handlerFunc != nil {
			nd.handlerFunc = n.handlerFunc
		}

		return nd
	}

	if nd.siblingNode == nil {
		n.parentNode = nd.parentNode
		nd.siblingNode = n
		return nd
	}

	nd.siblingNode = nd.siblingNode.combine(n)

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
