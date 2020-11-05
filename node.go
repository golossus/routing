package routing

import (
	"net/http"
	"reflect"
	"regexp"
)

const (
	catchAllExpression = "^.*$"
)

/*
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
*/
type nodeInterface interface {
	hasParameters() bool
	setParent(parent nodeInterface)
	sibling() nodeInterface
	setSibling(sibling nodeInterface)
	child() nodeInterface
	find(p string) nodeInterface
	combine(n nodeInterface) nodeInterface
	match(p string) bool
	handler() http.HandlerFunc
	setHandler(handler http.HandlerFunc)
	getPrefix() string
	getParent() nodeInterface
	setPrefix(prefix string)
	getWeight() int
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
	if ns.parentNode != nil {
		return ns.parentNode.hasParameters()
	}
	return false
}

func (ns *nodeStatic) setParent(parent nodeInterface) {
	ns.parentNode = parent
}

func (ns *nodeStatic) getPrefix() string {
	return ns.prefix
}

func (ns *nodeStatic) getParent() nodeInterface {
	return ns.parentNode
}

func (ns *nodeStatic) setPrefix(prefix string) {
	ns.prefix = prefix
}

func (ns *nodeStatic) sibling() nodeInterface {
	return ns.siblingNode
}

func (ns *nodeStatic) setSibling(sibling nodeInterface) {
	ns.siblingNode = sibling
}

func (ns *nodeStatic) setHandler(handlerFunc http.HandlerFunc) {
	ns.handlerFunc = handlerFunc
}

func (ns *nodeStatic) getWeight() int {
	return ns.weight
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
		if ns.sibling() == nil {
			return nil
		}

		return ns.sibling().find(p)
	}

	if pos == len(p) && len(p) == len(ns.prefix) {
		if nil != ns.handlerFunc {
			return ns
		}

		return nil
	}

	if ns.child() != nil {
		n := ns.child().find(p[pos:])
		if nil != n && nil != n.handler() {
			return n
		}
	}

	for next := ns.sibling(); nil != next; next = next.sibling() {
		if n := next.find(p); nil != n {
			return n
		}
	}

	return nil
}

func (ns *nodeStatic) combine(n nodeInterface) nodeInterface {

	if n == nil {
		return ns
	}

	if reflect.TypeOf(ns) != reflect.TypeOf(n) {

		if ns.sibling() != nil {
			ns.sibling().combine(n)
			ns.siblingNode.setParent(ns.parentNode)
			return ns
		}
		ns.siblingNode = n
		ns.siblingNode.setParent(ns.parentNode)
		return ns
	}

	pos := common(ns.prefix, n.getPrefix())

	if pos == 0 {
		if ns.sibling() != nil {
			ns.siblingNode = ns.sibling().combine(n)
			return ns
		}
		ns.siblingNode = n
		ns.siblingNode.setParent(ns.parentNode)
		return ns
	}

	if pos == len(ns.prefix) && pos != len(n.getPrefix()) {
		n.setPrefix(n.getPrefix()[pos:])
		n.setParent(ns)

		if ns.child() != nil {
			ns.childNode = ns.childNode.combine(n)
			return ns
		}
		ns.childNode = n
		return ns
	}

	if pos != len(ns.getPrefix()) && pos == len(n.getPrefix()) {
		ns.setPrefix(ns.getPrefix()[pos:])
		n.setSibling(ns.sibling())
		ns.setSibling(nil)
		ns.setParent(n)
		n.(*nodeStatic).childNode = ns.combine(n.child())

		return n
	}

	if pos != len(ns.getPrefix()) && pos != len(n.getPrefix()) {
		split := &nodeStatic{
			prefix:      ns.prefix[:pos],
			handlerFunc: nil,
			childNode:   nil,
			parentNode:  ns.parentNode,
			siblingNode: ns.siblingNode,
			weight:      0,
		}

		ns.prefix = ns.prefix[pos:]
		ns.parentNode = split
		ns.siblingNode = nil

		n.setPrefix(n.getPrefix()[pos:])
		n.setParent(split)

		split.childNode = ns.combine(n)

		return split
	}

	if n.handler() != nil {
		ns.handlerFunc = n.handler()
	}
	if ns.childNode != nil {
		ns.childNode.combine(n.child())
		ns.child().setParent(ns)
		return ns
	}
	ns.childNode = n.child()
	return ns
}

func (ns *nodeStatic) match(p string) bool {
	return false
}

func common(s1, s2 string) int {
	for k := 0; k < len(s1); k++ {
		if k == len(s2) || s1[k] != s2[k] {
			return k
		}
	}

	return len(s1)
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

func (nd *nodeDynamic) setPrefix(prefix string) {
	nd.prefix = prefix
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
			if nd.sibling() == nil {
				return nil
			}
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

	if nd.sibling() == nil {
		return nil
	}
	return nd.sibling().find(p)
}

func (nd *nodeDynamic) combine(n nodeInterface) nodeInterface {
	if n == nil {
		return nd
	}

	if reflect.TypeOf(nd) == reflect.TypeOf(n) {
		if nd.getPrefix() == n.getPrefix() {
			if !nd.regexpEquals(n.(*nodeDynamic)) {
				if nd.siblingNode == nil {
					nd.siblingNode = n
				} else {
					nd.siblingNode = nd.siblingNode.combine(n)
				}
				nd.siblingNode.setParent(nd)
				return nd
			}

			for k := range n.(*nodeDynamic).childrenNodes {
				n.(*nodeDynamic).childrenNodes[k].setParent(nd)
			}

			for k, next1 := range nd.childrenNodes {
				if next2, ok := n.(*nodeDynamic).childrenNodes[k]; !ok {
					n.(*nodeDynamic).childrenNodes[k] = next1
				} else {
					n.(*nodeDynamic).childrenNodes[k] = next1.combine(next2)
				}
			}

			nd.childrenNodes = n.(*nodeDynamic).childrenNodes
			if n.handler() != nil {
				nd.handlerFunc = n.handler()
			}

			return nd
		}

		if nd.siblingNode == nil {
			nd.siblingNode = n
		} else {
			nd.siblingNode = nd.siblingNode.combine(n)
		}
		nd.siblingNode.setParent(n.(*nodeDynamic).parentNode)
		return nd
	}

	n.(*nodeStatic).siblingNode = nd
	n.setParent(nd.parentNode)
	return n

}

func (nd *nodeDynamic) regexpEquals(o *nodeDynamic) bool {
	return nd.regexpToString() == o.regexpToString()
}

func (nd *nodeDynamic) regexpToString() string {
	if nd.regexp == nil {
		return ""
	}
	return nd.regexp.String()
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
