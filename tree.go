package hw14_go

const (
	NodeTypeStatic = iota
	NodeTypeDynamic
)

type Tree struct {
	root map[string]*Node
}

func (t *Tree) Insert(verb string, chunks []chunk, handler HandlerFunction) {

	if nil == t.root {
		t.root = make(map[string]*Node)
	}

	n, _ := t.root[verb]

	var leaf *Node
	h := handler
	if len(chunks) > 1 {
		h = nil
	}

	t.root[verb], leaf = insert(n, chunks[0].v, h)
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

func (t *Tree) Find(verb string, path string) (HandlerFunction, urlParameterBag) {
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
					urlParameter := urlParameter{name: n.prefix, value: p[0:i]}
					urlParameterBag.addParameter(urlParameter)
					n = next
					p = p[i:]
					found = true
					break
				}
			}
			if !found {
				urlParameter := urlParameter{name: n.prefix, value: p}
				urlParameterBag.addParameter(urlParameter)
				return n.handler, urlParameterBag
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
}

func insert(n *Node, path string, handler HandlerFunction) (root, leaf *Node) {

	if nil == n {
		leaf = &Node{prefix: path, handler: handler, t: NodeTypeStatic}
		return leaf, leaf
	}

	if NodeTypeDynamic == n.t {
		if n.child == nil {
			n.child = &Node{prefix: path, t: NodeTypeStatic, handler: handler}
			leaf = n.child
			n.stops[path[0:1]] = leaf
			return n, leaf
		}
		n.child, leaf = insert(n.child, path, handler)
		n.stops[path[0:1]] = leaf
		return n, leaf
	}

	pos := common(n.prefix, path)

	if pos == 0 {
		if n.sibling != nil && n.sibling.t == NodeTypeDynamic {
			n.sibling, leaf = insertSibling(n.sibling, path, handler)
			return n, leaf
		}
		n.sibling, leaf = insert(n.sibling, path, handler)
		return n, leaf
	}

	if pos < len(n.prefix) {
		newNode := &Node{prefix: n.prefix[0:pos], child: n, t: NodeTypeStatic, sibling: n.sibling}
		n.prefix = n.prefix[pos:]
		n.sibling = nil
		n = newNode
	}

	if pos == len(path) {
		if n.handler == nil {
			n.handler = handler
		}
		return n, n
	}

	n.child, leaf = insert(n.child, path[pos:], handler)

	return n, leaf
}

func insertSibling(sibling *Node, path string, handler HandlerFunction) (root, leaf *Node) {
	if sibling.sibling == nil {
		sibling.sibling = &Node{prefix: path, t: NodeTypeStatic, handler: handler}
		leaf = sibling.sibling
		return sibling, leaf
	}

	if sibling.sibling.t == NodeTypeDynamic {
		sibling.sibling, leaf = insertSibling(sibling.sibling, path, handler)
		return sibling, leaf
	}

	sibling.sibling, leaf = insert(sibling.sibling, path, handler)
	return sibling, leaf
}

func insertDynamic(n *Node, ident string, handler HandlerFunction) (root, leaf *Node) {

	if n.child == nil {
		stops := make(map[string]*Node)
		n.child = &Node{prefix: ident, t: NodeTypeDynamic, handler: handler, stops: stops}
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
			stops := make(map[string]*Node)
			tmp.sibling = &Node{prefix: ident, t: NodeTypeDynamic, handler: handler, stops: stops}
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
