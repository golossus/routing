package routing

import (
	"net/http"
)

type tree struct {
	root nodeInterface
}
// [Santi][Concern] I don't know if this the best signature. Maybe, as find method, this
// should receive the path and hadler to insert instead of an already create tree root
func (t *tree) insert(root nodeInterface) nodeInterface {

	if nil == t.root {
		t.root = root
		return t.root
	}

	t.root = t.root.merge(root)
	return t.root
}

func (t *tree) find(path string) nodeInterface {
	return t.root.find(path)
}

func createTreeFromPath(path string, handler http.HandlerFunc) (root, leaf nodeInterface, err error) {

	parser := newParser(path)
	_, err = parser.parse()
	if err != nil {
		return nil, nil, err
	}

	var h http.HandlerFunc

	c := parser.chunks
	for i := 0; i < len(c); i++ {

		if len(c)-1 == i {
			h = handler
		}

		if i == 0 {
			root = createNodeFromChunk(c[i], h)
			leaf = root
			continue
		}

		if nil != leaf {
			next := createNodeFromChunk(c[i], h)
			leaf.addChild(next)
			leaf = next
		}
	}

	return root, leaf, nil
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
