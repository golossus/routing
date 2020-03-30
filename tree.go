package routing

import (
	"fmt"
	"net/http"
	"regexp"
)

const (
	nodeTypeStatic = iota
	nodeTypeDynamic
)

type tree struct {
	root *node
}

func (t *tree) insert(chunks []chunk, handler http.HandlerFunc) {
	subtree, err := createTreeFromChunks(chunks, handler)
	if err != nil {
		panic(err)
	}
	t.root = combine(t.root, subtree)
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
			for k, v := range tree1.stops {
				tree2.stops[k] = v
			}
			tree1.stops = tree2.stops
			tree1.child = combine(tree1.child, tree2.child)
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

	tree1.child = combine(tree1.child, tree2.child)
	return tree1
}

func createTreeFromChunks(chunks []chunk, handler http.HandlerFunc) (*node, error) {

	if len(chunks) < 1 {
		return nil, fmt.Errorf("chunks can not be empty")
	}

	var root = createNodeFromChunk(chunks[0])
	n := root

	for i := 1; i < len(chunks); i++ {
		newNode := createNodeFromChunk(chunks[i])
		if n.t == nodeTypeDynamic {
			n.stops[newNode.prefix[0]] = newNode
		}
		n.child = newNode
		n = n.child
	}

	n.handler = handler

	return root, nil
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

type node struct {
	prefix  string
	handler http.HandlerFunc
	child   *node
	sibling *node
	t       int
	stops   map[byte]*node
	regexp  *regexp.Regexp
}

func common(s1, s2 string) int {
	for k := 0; k < len(s1); k++ {
		if k == len(s2) || s1[k] != s2[k] {
			return k
		}
	}

	return len(s1)
}
