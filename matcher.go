package routing

import "net/http"

type matcher func(r *http.Request) bool

func byHost(host string) (matcher, error) {
	parser := newParser("/"+host)
	_, err := parser.parse()
	if err != nil {
		return nil, err
	}

	root, leaf := createTreeFromChunks(parser.chunks)
	leaf.handler = func(writer http.ResponseWriter, request *http.Request) {}

	return func(r *http.Request) bool {
		return nil != find(root, "/"+r.Host, r)
	}, nil
}
