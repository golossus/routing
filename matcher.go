package routing

import (
	"net/http"
	"strings"
)

type matcher func(r *http.Request) (bool, *node)

// CustomMatcher defines the type of a custom function to match against an http
// request and allow a route to be found or not in the Router
type CustomMatcher func(r *http.Request) bool

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

func byHeaders(headers map[string]string) matcher {

	return func(r *http.Request) (bool, *node) {
		if r == nil || len(headers) > len(r.Header) {
			return false, nil
		}
		for key, value := range headers {
			if r.Header.Get(key) != value {
				return false, nil
			}
		}
		return true, nil
	}
}

func byQueryParameters(params map[string]string) matcher {

	return func(r *http.Request) (bool, *node) {
		if r == nil || len(params) > len(r.URL.Query()) {
			return false, nil
		}
		for key, value := range params {
			if r.URL.Query().Get(key) != value {
				return false, nil
			}
		}
		return true, nil
	}
}

func byCustomMatcher(custom func(r *http.Request) bool) matcher {

	return func(r *http.Request) (bool, *node) {
		if r == nil {
			return false, nil
		}

		return custom(r), nil
	}
}
