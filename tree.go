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

func (t *tree) insert(chunks []chunk, handler http.HandlerFunc) {
	t.root = combine(t.root, createTreeFromChunks(chunks, handler))
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
				return tree1
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
			return tree1
		}

		if tree2.t == nodeTypeStatic {
			tree2.sibling = tree1
			return tree2
		}
	}

	if tree2.t == nodeTypeDynamic {
		tree1.sibling = combine(tree1.sibling, tree2)
		return tree1
	}

	pos := common(tree1.prefix, tree2.prefix)

	if pos == 0 {
		tree1.sibling = combine(tree1.sibling, tree2)
		return tree1
	}

	if pos == len(tree1.prefix) && pos != len(tree2.prefix) {
		tree2.prefix = tree2.prefix[pos:]
		tree1.child = combine(tree1.child, tree2)
		return tree1
	}

	if pos != len(tree1.prefix) && pos == len(tree2.prefix) {
		tree1.prefix = tree1.prefix[pos:]
		tree2.sibling = tree1.sibling
		tree1.sibling = nil
		tree2.child = combine(tree1, tree2.child)
		return tree2
	}

	if pos != len(tree1.prefix) && pos != len(tree2.prefix) {
		split := createNodeFromChunk(chunk{t: tChunkStatic, v: tree1.prefix[:pos]})
		split.sibling = tree1.sibling

		tree1.prefix = tree1.prefix[pos:]
		tree1.sibling = nil

		tree2.prefix = tree2.prefix[pos:]

		split.child = combine(tree1, tree2)
		return split
	}

	if tree2.handler != nil {
		tree1.handler = tree2.handler
	}

	tree1.child = combine(tree1.child, tree2.child)
	return tree1
}

func createTreeFromChunks(chunks []chunk, handler http.HandlerFunc) *node {

	if len(chunks) < 1 {
		return nil
	}

	var root = createNodeFromChunk(chunks[0])
	n := root

	for i := 1; i < len(chunks); i++ {
		newNode := createNodeFromChunk(chunks[i])
		if n.t == nodeTypeDynamic {
			n.stops[newNode.prefix[0]] = newNode
		} else {
			n.child = newNode
		}

		n = newNode
	}

	n.handler = handler

	return root
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

func (t *tree) find(path string) (http.HandlerFunc, URLParameterBag) {
	urlParameterBag := newURLParameterBag(5, true)

	return find(t.root, path, &urlParameterBag), urlParameterBag
}

func find(n *node, p string, urlParameterBag *URLParameterBag) http.HandlerFunc {
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
					h := find(next, p[i:], urlParameterBag)
					if nil != h {
						urlParameterBag.add(n.prefix, p[0:i])
						return h
					}
				}
			}

			if p[i] == '/' && !n.isCatchAll() {
				return find(n.sibling, p, urlParameterBag)
			}
		}

		if nil != n.handler && !traversed {
			validExpression := true
			if n.regexp != nil {
				validExpression = n.regexp.MatchString(p)
			}
			if validExpression {
				urlParameterBag.add(n.prefix, p)
				return n.handler
			}
		}

		return find(n.sibling, p, urlParameterBag)
	}

	pos := common(p, n.prefix)
	if pos == 0 {
		return find(n.sibling, p, urlParameterBag)
	}

	if pos == len(p) && len(p) == len(n.prefix) {
		return n.handler
	}

	h := find(n.child, p[pos:], urlParameterBag)
	if nil != h {
		return h
	}

	for next := n.sibling; nil != next; next = next.sibling {
		if next.t != nodeTypeDynamic {
			continue
		}

		return find(next, p, urlParameterBag)
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
