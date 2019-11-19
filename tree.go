package hw14_go

type Tree struct {
	root *Node
}

func (t *Tree) Insert(path string, handler HandlerFunction) {
	t.root = insert(t.root, path, handler)
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
}

func insert(n *Node, path string, handler HandlerFunction) *Node {

	if nil == n {
		return &Node{prefix: path, handler: handler}
	}

	pos := common(n.prefix, path)

	if pos == len(path) {
		n.handler = handler
		return n
	}

	if pos == 0 {
		n.sibling = insert(n.sibling, path, handler)
		return n
	}

	if pos < len(n.prefix) {
		newNode := &Node{prefix: n.prefix[0:pos], child: n}
		n.prefix = n.prefix[pos:]
		n = newNode

	}

	n.child = insert(n.child, path[pos:], handler)

	return n
}

func common(s1, s2 string) int {
	for k := 0; k < len(s1); k++ {
		if k == len(s2) || s1[k] != s2[k] {
			return k
		}
	}

	return len(s1)
}
