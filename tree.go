package hw14_go

import (
	"fmt"
)

type Node struct {
	prefix  string
	handler HandlerFunction
	child   *Node
	sibling *Node
}

func (n *Node) insert(path string, handler HandlerFunction) *Node {

		if nil == n {
			n = &Node{prefix: path, handler: handler}
			return n
		}

		pos := common(n.prefix, path)

		if pos == 0 {
			if nil == node.sibling {
				node.sibling = &Node{path: path, handler: handler}
				return nil
			}

			node = node.sibling;
			continue
		}

		if pos < len(node.path) {
			n := &Node{path: node.path[0:pos], child: node}
			node.path = node.path[pos:]
			node = n

			return nil
		}

		if pos == len(path) {
			if nil != node.handler {
				return fmt.Errorf("handler already defined for %s", path)
			}
			node.handler = handler
		}

		path = path[pos:]
		node = node.sibling


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
