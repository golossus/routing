package routing

import (
	"net/http"
)

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

	pos := ns.common(p)

	if pos == 0 {
		if ns.siblingNode == nil {
			return nil
		}

		return ns.siblingNode.find(p)
	}

	if pos == len(p) && len(p) == len(ns.prefix) {
		if nil != ns.handlerFunc {
			return ns
		}

		return nil
	}

	if ns.childNode != nil {
		n := ns.childNode.find(p[pos:])
		if nil != n  {
			return n
		}
	}

	if ns.siblingNode != nil {
		return ns.siblingNode.find(p)
	}

	return nil
}

func (ns *nodeStatic) combine(ni nodeInterface) nodeInterface {

	if ni == nil {
		return ns
	}

	n, ok := ni.(*nodeStatic)
	if !ok {
		if ns.siblingNode != nil {
			ns.siblingNode.combine(ni)
			return ns
		}

		ns.siblingNode = ni
		ns.siblingNode.setParent(ns.parentNode)
		return ns
	}

	pos := ns.common(n.prefix)

	if pos == 0 {
		if ns.siblingNode != nil {
			ns.siblingNode = ns.siblingNode.combine(n)
			return ns
		}
		n.parentNode = ns.parentNode
		ns.siblingNode = n
		return ns
	}

	if pos == len(ns.prefix) && pos != len(n.prefix) {
		n.prefix = n.prefix[pos:]
		n.parentNode = ns

		if ns.childNode != nil {
			ns.childNode = ns.childNode.combine(n)
			return ns
		}
		ns.childNode = n
		return ns
	}

	if pos != len(ns.prefix) && pos == len(n.prefix) {
		ns.prefix = ns.prefix[pos:]
		n.siblingNode = ns.siblingNode
		ns.siblingNode = nil
		ns.parentNode = n
		n.childNode = ns.combine(n.childNode)

		return n
	}

	if pos != len(ns.prefix) && pos != len(n.prefix) {
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

		n.prefix = n.prefix[pos:]
		n.parentNode = split

		split.childNode = ns.combine(n)

		return split
	}

	if n.handlerFunc != nil {
		ns.handlerFunc = n.handlerFunc
	}

	if ns.childNode != nil {
		ns.childNode.combine(n.childNode)
		ns.childNode.setParent(ns)
		return ns
	}

	ns.childNode = n.childNode
	return ns
}

func (ns *nodeStatic) common(p string) int {
	for k := 0; k < len(ns.prefix); k++ {
		if k == len(p) || ns.prefix[k] != p[k] {
			return k
		}
	}

	return len(ns.prefix)
}
