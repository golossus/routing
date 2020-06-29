package routing

import (
	"net/http"
)

const (
	nodeTypeStatic = iota
	nodeTypeDynamic
)

type tree struct {
	root *node
}

func (t *tree) insert(chunks []chunk, handler http.HandlerFunc) *node {
	root2, leaf2 := createTreeFromChunks(chunks)
	leaf2.handler = handler

	t.root = combine(t.root, root2)

	return leaf2
}

func combine(tree1 *node, tree2 *node) *node {

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

func createTreeFromChunks(chunks []chunk) (root, leaf *node) {

	if len(chunks) < 1 {
		return nil, nil
	}

	root = createNodeFromChunk(chunks[0])
	n := root

	for i := 1; i < len(chunks); i++ {
		newNode := createNodeFromChunk(chunks[i])
		if n.t == nodeTypeDynamic {
			n.stops[newNode.prefix[0]] = newNode
		} else {
			n.child = newNode
		}
		newNode.parent = n
		n = newNode
	}

	return root, n
}

func createNodeFromChunk(c chunk) *node {
	var n *node
	if c.t == tChunkStatic {
		n = &node{prefix: c.v, handler: nil, t: nodeTypeStatic}
	} else {
		stops := make(map[byte]*node)

		n = &node{prefix: c.v, t: nodeTypeDynamic, handler: nil, stops: stops, regexp: c.exp}
	}
	return n
}

func (t *tree) find(path string) *node {
	return find(t.root, path)
}

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

func common(s1, s2 string) int {
	for k := 0; k < len(s1); k++ {
		if k == len(s2) || s1[k] != s2[k] {
			return k
		}
	}

	return len(s1)
}

func calcWeight(n *node) int {
	if n == nil {
		return 0
	}

	n.w = 0
	if n.handler != nil {
		n.w++
	}

	if n.t == nodeTypeStatic {
		n.w = n.w + calcWeight(n.child) + calcSiblingsWeight(n.child)
	} else {
		for _, c := range n.stops {
			n.w = n.w + calcWeight(c) + calcSiblingsWeight(c)
		}
	}

	return n.w
}

func calcSiblingsWeight(n *node) int {
	if n == nil {
		return 0
	}

	w := 0
	s := n.sibling
	for s != nil {
		if s.t == nodeTypeStatic {
			w = w + calcWeight(s)
		} else {
			for _, c := range s.stops {
				w = w + calcWeight(c)
			}
		}

		s = s.sibling
	}

	return w
}

func sortByWeight(head *node) *node {
	var sorted *node

	current := head
	for current != nil {
		next := current.sibling

		if current.t == nodeTypeStatic {
			current.child = sortByWeight(current.child)
		} else {
			for k, s := range current.stops {
				current.stops[k] = sortByWeight(s)
			}
		}
		sorted = sortInsertByWeight(sorted, current)

		current = next
	}

	return sorted
}

func sortInsertByWeight(head *node, in *node) *node {
	var current *node
	if head == nil || head.w < in.w {
		in.sibling = head
		head = in
	} else {
		current = head
		for current.sibling != nil && current.sibling.w >= in.w {
			current = current.sibling
		}
		in.sibling = current.sibling
		current.sibling = in
	}

	return head
}
