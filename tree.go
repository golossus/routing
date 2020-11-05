package routing

import (
	"net/http"
)

type tree struct {
	root nodeInterface
}

func (t *tree) insert(chunks []chunk, handler http.HandlerFunc) nodeInterface {
	root2, leaf2 := createTreeFromChunks(chunks)
	leaf2.setHandler(handler)

	if t.root != nil {
		t.root = t.root.combine(root2)
	} else {
		t.root = root2
	}

	return leaf2
}

/*
func combine(tree1 nodeInterface, tree2 nodeInterface) nodeInterface {

	return tree1.combine(tree2)

	if tree1 == nil {
		return tree2
	}

	if tree2 == nil {
		return tree1
	}

	if tree1.t == nodeTypeDynamic {
		if tree2.t == nodeTypeDynamic && tree2.prefix == tree1.prefix {
			if !tree1.regexpEquals(tree2) {
				tree1.sibling = combine(tree1.sibling, tree2)
				tree1.sibling.parent = tree1
				return tree1
			}

			for k := range tree2.stops {
				tree2.stops[k].parent = tree1
			}

			for k, next1 := range tree1.stops {
				if next2, ok := tree2.stops[k]; !ok {
					tree2.stops[k] = next1
				} else {
					tree2.stops[k] = combine(next1, next2)
				}
			}

			tree1.stops = tree2.stops
			if tree2.handler != nil {
				tree1.handler = tree2.handler
			}

			return tree1
		}

		if tree2.t == nodeTypeDynamic && tree2.prefix != tree1.prefix {
			tree1.sibling = combine(tree1.sibling, tree2)
			tree1.sibling.parent = tree1.parent
			return tree1
		}

		if tree2.t == nodeTypeStatic {
			tree2.sibling = tree1
			tree2.parent = tree1.parent
			return tree2
		}
	}

	if tree2.t == nodeTypeDynamic {
		tree1.sibling = combine(tree1.sibling, tree2)
		tree1.sibling.parent = tree1.parent
		return tree1
	}

	pos := common(tree1.prefix, tree2.prefix)

	if pos == 0 {
		tree1.sibling = combine(tree1.sibling, tree2)
		tree1.sibling.parent = tree1.parent
		return tree1
	}

	if pos == len(tree1.prefix) && pos != len(tree2.prefix) {
		tree2.prefix = tree2.prefix[pos:]
		tree2.parent = tree1
		tree1.child = combine(tree1.child, tree2)
		return tree1
	}

	if pos != len(tree1.prefix) && pos == len(tree2.prefix) {
		tree1.prefix = tree1.prefix[pos:]
		tree2.sibling = tree1.sibling
		tree1.sibling = nil
		tree1.parent = tree2
		tree2.child = combine(tree1, tree2.child)
		return tree2
	}
	if pos != len(tree1.prefix) && pos != len(tree2.prefix) {
		split := createNodeFromChunk(chunk{t: tChunkStatic, v: tree1.prefix[:pos]})
		split.parent = tree1.parent
		split.sibling = tree1.sibling

		tree1.prefix = tree1.prefix[pos:]
		tree1.parent = split
		tree1.sibling = nil

		tree2.prefix = tree2.prefix[pos:]
		tree2.parent = split

		split.child = combine(tree1, tree2)

		return split
	}

	if tree2.handler != nil {
		tree1.handler = tree2.handler
	}

	tree1.child = combine(tree1.child, tree2.child)
	tree1.child.parent = tree1
	return tree1
}
*/
func createTreeFromChunks(chunks []chunk) (root, leaf nodeInterface) {

	if len(chunks) < 1 {
		return nil, nil
	}

	root = createNodeFromChunk(chunks[0])
	n := root

	for i := 1; i < len(chunks); i++ {
		newNode := createNodeFromChunk(chunks[i])
		switch n.(type) {
		case *nodeDynamic:
			n.(*nodeDynamic).childrenNodes[newNode.getPrefix()[0]] = newNode
		case *nodeStatic:
			n.(*nodeStatic).childNode = newNode
		}
		newNode.setParent(n)
		n = newNode
	}

	return root, n
}

func createNodeFromChunk(c chunk) nodeInterface {
	var n nodeInterface
	if c.t == tChunkStatic {
		n = &nodeStatic{
			prefix:      c.v,
			handlerFunc: nil,
			childNode:   nil,
			parentNode:  nil,
			siblingNode: nil,
			weight:      0,
		}
	} else {
		stops := make(map[byte]nodeInterface)
		n = &nodeDynamic{
			prefix:        c.v,
			handlerFunc:   nil,
			parentNode:    nil,
			siblingNode:   nil,
			childrenNodes: stops,
			regexp:        c.exp,
			weight:        0,
		}
	}
	return n
}

func (t *tree) find(path string) nodeInterface {
	return t.root.find(path)
}

/*
func find(n *node, p string) *node {
	if nil == n || len(p) == 0 {
		return nil
	}

	if n.t == nodeTypeDynamic {
		traversed := false
		for i := 0; i < len(p); i++ {

			if next, ok := n.stops[p[i]]; ok {
				validExpression := true
				if n.regexp != nil {
					validExpression = n.regexp.MatchString(p[0:i])
				}
				if validExpression {
					traversed = true
					h := find(next, p[i:])
					if nil != h && nil != h.handler {
						return h
					}
				}
			}

			if p[i] == '/' && !n.isCatchAll() {
				return find(n.sibling, p)
			}
		}

		if nil != n.handler && !traversed {
			validExpression := true
			if n.regexp != nil {
				validExpression = n.regexp.MatchString(p)
			}
			if validExpression {
				return n
			}
		}

		return find(n.sibling, p)
	}

	pos := common(p, n.prefix)
	if pos == 0 {
		return find(n.sibling, p)
	}

	if pos == len(p) && len(p) == len(n.prefix) {
		if nil != n.handler {
			return n
		}

		return nil
	}

	h := find(n.child, p[pos:])
	if nil != h && nil != h.handler {
		return h
	}

	for next := n.sibling; nil != next; next = next.sibling {
		if next.t != nodeTypeDynamic {
			continue
		}

		return find(next, p)
	}

	return nil
}
*/

func calcWeight(n nodeInterface) int {
	if n == nil {
		return 0
	}

	switch n.(type) {
	case *nodeStatic:
		n.(*nodeStatic).weight = 0
		if n.handler() != nil {
			n.(*nodeStatic).weight++
		}
		n.(*nodeStatic).weight = n.(*nodeStatic).weight + calcWeight(n.child()) + calcSiblingsWeight(n.child())
		return n.(*nodeStatic).weight

	case *nodeDynamic:
		n.(*nodeDynamic).weight = 0
		if n.handler() != nil {
			n.(*nodeDynamic).weight++
		}
		for _, c := range n.(*nodeDynamic).childrenNodes {
			n.(*nodeDynamic).weight = n.(*nodeDynamic).weight + calcWeight(c) + calcSiblingsWeight(c)
		}
		return n.(*nodeDynamic).weight
	}

	return 0
}

func calcSiblingsWeight(n nodeInterface) int {
	if n == nil {
		return 0
	}

	w := 0
	s := n.sibling()
	for s != nil {
		switch n.(type) {
		case *nodeStatic:
			w = w + calcWeight(s)
		case *nodeDynamic:
			for _, c := range s.(*nodeDynamic).childrenNodes {
				w = w + calcWeight(c)
			}
		}

		s = s.sibling()
	}

	return w
}

func sortByWeight(head nodeInterface) nodeInterface {
	var sorted nodeInterface

	current := head
	for current != nil {
		next := current.sibling()

		switch current.(type) {
		case *nodeStatic:
			current.(*nodeStatic).childNode = sortByWeight(current.child())
		case *nodeDynamic:
			for k, s := range current.(*nodeDynamic).childrenNodes {
				current.(*nodeDynamic).childrenNodes[k] = sortByWeight(s)
			}
		}

		sorted = sortInsertByWeight(sorted, current)

		current = next
	}

	return sorted
}

func sortInsertByWeight(head nodeInterface, in nodeInterface) nodeInterface {
	var current nodeInterface

	if head == nil || head.getWeight() < in.getWeight() {
		in.setSibling(head)
		head = in
	} else {
		current = head
		for current.sibling() != nil && current.sibling().getWeight() >= in.getWeight() {
			current = current.sibling()
		}
		in.setSibling(current.sibling())
		current.setSibling(in)
	}

	return head
}
