package routing

import (
	"net/http"
)

type nodeInterface interface {
	find(p string) nodeInterface
	combine(n nodeInterface) nodeInterface
	hasParameters() bool
	setParent(parent nodeInterface)
	sibling() nodeInterface
	setSibling(sibling nodeInterface)
	child() nodeInterface
	handler() http.HandlerFunc
	setHandler(handler http.HandlerFunc)
	getPrefix() string
	getParent() nodeInterface
	getWeight() int
}
