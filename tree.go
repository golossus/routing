package hw14_go

const (
	NodeTypeStatic = iota
	NodeTypeDynamic
)

type Tree struct {
	root *Node
}

func (t *Tree) Insert(chunks []chunk, handler HandlerFunction) {
	var leaf *Node
	h := handler
	if len(chunks) > 1 {
		h = nil
	}

	t.root, leaf = insert(t.root, chunks[0].v, h)
	chunks = chunks[1:]
	for index, chunk := range chunks {
		if index == len(chunks)-1 {
			h = handler
		}
		next := leaf

		if chunk.t == TChunkStatic {
			_, leaf = insert(next, chunk.v, h)
			continue
		}

		_, leaf = insertDynamic(next, chunk.v, h)
	}

}

func (t *Tree) Find(path string) HandlerFunction {
	n := t.root
	p := path

	for nil != n && len(p) > 0 {
		pos := common(p, n.prefix)
		if pos == 0 {
			n = n.sibling
			continue
		}

		if pos == len(p) {
			return n.handler
		}

		p = p[pos:]
		n = n.child
	}

	return nil
}

type Node struct {
	prefix  string
	handler HandlerFunction
	child   *Node
	sibling *Node
	t       int
}

func insert(n *Node, path string, handler HandlerFunction) (root, leaf *Node) {

	if nil == n {
		leaf = &Node{prefix: path, handler: handler, t: NodeTypeStatic}
		return leaf, leaf
	}

	pos := common(n.prefix, path)

	if pos == len(path) {
		n.handler = handler
		return n, n
	}

	if pos == 0 {
		n.sibling, leaf = insert(n.sibling, path, handler)
		return n, leaf
	}

	if pos < len(n.prefix) {
		newNode := &Node{prefix: n.prefix[0:pos], child: n, t: NodeTypeStatic}
		n.prefix = n.prefix[pos:]
		n = newNode
	}

	n.child, leaf = insert(n.child, path[pos:], handler)

	return n, leaf
}

func insertDynamic(n *Node, ident string, handler HandlerFunction) (root, leaf *Node) {

	if n.child == nil {
		n.child = &Node{prefix: ident, t: NodeTypeDynamic, handler: handler}
		leaf = n.child
	}

	tmp := n.child

	for {
		if tmp.t == NodeTypeDynamic && tmp.prefix == ident {
			if tmp.handler == nil {
				tmp.handler = handler
			}
			leaf = tmp
			return n, leaf
		}

		if tmp.sibling == nil {
			tmp.sibling = &Node{prefix: ident, t: NodeTypeDynamic, handler: handler}
			leaf = tmp.sibling
			return n, leaf
		}

		tmp = tmp.sibling
	}
}

func common(s1, s2 string) int {
	for k := 0; k < len(s1); k++ {
		if k == len(s2) || s1[k] != s2[k] {
			return k
		}
	}

	return len(s1)
}
