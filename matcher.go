package routing

import (
	"net/http"
	"strings"
)

type matcher func(r *http.Request) (bool, *node)

func byHost(host string) (matcher, error) {
	parser := newParser("/" + strings.ToLower(host))
	_, err := parser.parse()
	if err != nil {
		return nil, err
	}

	root, leaf := createTreeFromChunks(parser.chunks)
	leaf.handler = func(writer http.ResponseWriter, request *http.Request) {}

	return func(r *http.Request) (bool, *node) {
		if r == nil {
			return false, leaf
		}

		return nil != find(root, "/"+r.Host, r), leaf
	}, nil
}

func bySchemas(schemas ...string) (matcher, error) {

	t := &tree{}

	for _, schema := range schemas {
		parser := newParser("/" + strings.ToLower(schema))
		_, err := parser.parse()
		if err != nil {
			return nil, err
		}
		t.insert(parser.chunks, func(writer http.ResponseWriter, request *http.Request) {})
	}

	return func(r *http.Request) (bool, *node) {
		if r == nil {
			return false, t.root
		}

		leaf := find(t.root, "/"+r.URL.Scheme, r)
		return nil != leaf, leaf
	}, nil
}
