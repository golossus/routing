package http_router

import (
	"fmt"
	"regexp"
)

const (
	NodeTypeStatic = iota
	NodeTypeDynamic
)

type Tree struct {
	root map[string]*Node
}

func (t *Tree) Insert(verb string, chunks []chunk, handler HandlerFunction) {

	subtree, err := createTreeFromChunks(chunks, handler)
	if err != nil {
		panic(err)
	}

	if nil == t.root {
		t.root = make(map[string]*Node)
	}

	t.root[verb] = combine(t.root[verb], subtree)
}

func combine(tree1 *Node, tree2 *Node) *Node {

	if tree1 == nil {
		return tree2
	}

	if tree2 == nil {
		return tree1
	}

	if tree1.t == NodeTypeDynamic {
		if tree2.t == NodeTypeDynamic && tree2.prefix == tree1.prefix {
			for k, v := range tree1.stops {
				tree2.stops[k] = v
			}
			tree1.stops = tree2.stops
			tree1.child = combine(tree1.child, tree2.child)
			return tree1
		}

		if tree2.t == NodeTypeDynamic && tree2.prefix != tree1.prefix {
			tree1.sibling = combine(tree1.sibling, tree2)
			return tree1
		}

		if tree2.t == NodeTypeStatic {
			tree2.sibling = tree1
			return tree2
		}
	}

	if tree2.t == NodeTypeDynamic {
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
		split := createNodeFromChunk(chunk{t: TChunkStatic, v: tree1.prefix[:pos]})
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

func createTreeFromChunks(chunks []chunk, handler HandlerFunction) (*Node, error) {

	if len(chunks) < 1 {
		return nil, fmt.Errorf("chunks can not be empty")
	}

	var root = createNodeFromChunk(chunks[0])
	n := root

	for i := 1; i < len(chunks); i++ {
		newNode := createNodeFromChunk(chunks[i])
		if n.t == NodeTypeDynamic {
			n.stops[newNode.prefix[0:1]] = newNode
		}
		n.child = newNode
		n = n.child
	}

	n.handler = handler

	return root, nil
}

func createNodeFromChunk(c chunk) *Node {
	var n *Node
	if c.t == TChunkStatic {
		n = &Node{prefix: c.v, handler: nil, t: NodeTypeStatic}
	} else {
		stops := make(map[string]*Node)

		n = &Node{prefix: c.v, t: NodeTypeDynamic, handler: nil, stops: stops, regexp: c.exp}
	}
	return n
}

func (t *Tree) Find(verb string, path string) (HandlerFunction, UrlParameterBag) {
	urlParameterBag := NewUrlParameterBag()

	n, ok := t.root[verb]
	if !ok {
		return nil, urlParameterBag
	}
	p := path

	for nil != n && len(p) > 0 {

		if n.t == NodeTypeDynamic {
			found := false
			for i, ch := range p {
				if next, ok := n.stops[string(ch)]; ok {
					validExpression := true
					if n.regexp != nil {
						validExpression = n.regexp.MatchString(p[0:i])
					}
					if validExpression {
						urlParameter := urlParameter{name: n.prefix, value: p[0:i]}
						urlParameterBag.addParameter(urlParameter)
						n = next
						p = p[i:]
						found = true
						break
					}
				}
			}
			if !found {
				validExpression := true
				if n.regexp != nil {
					validExpression = n.regexp.MatchString(p)
				}
				if validExpression {
					urlParameter := urlParameter{name: n.prefix, value: p}
					urlParameterBag.addParameter(urlParameter)
					return n.handler, urlParameterBag
				}
				return nil, urlParameterBag
			}
		}

		pos := common(p, n.prefix)
		if pos == 0 {
			n = n.sibling
			continue
		}

		if pos == len(p) && len(p) == len(n.prefix) {
			return n.handler, urlParameterBag
		}

		p = p[pos:]
		n = n.child
	}

	return nil, urlParameterBag
}

type Node struct {
	prefix  string
	handler HandlerFunction
	child   *Node
	sibling *Node
	t       int
	stops   map[string]*Node
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
