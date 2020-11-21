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
		t.root = t.root.merge(root2)
	} else {
		t.root = root2
	}

	return leaf2
}

func createTreeFromChunks(chunks []chunk) (root, leaf nodeInterface) {

	if len(chunks) < 1 {
		return nil, nil
	}

	for i := 0; i < len(chunks); i++ {
		if i == 0 {
			root = createNodeFromChunk(chunks[i])
			leaf = root
			continue
		}

		_, leaf = leaf.addChild(createNodeFromChunk(chunks[i]))
	}

	return root, leaf

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
