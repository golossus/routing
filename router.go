package hw14_go

import (
	"net/http"
)

type HandlerFunction func(http.ResponseWriter, *http.Request)

type Router interface {
	http.Handler
	AddHandler(string, HandlerFunction)
}

type routeHandler struct {
	path    string
	handler HandlerFunction
}

type SliceRouter struct {
	handlers []*routeHandler
}

func (r *SliceRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	for i := 0; i < len(r.handlers); i++ {
		if r.handlers[i].path == request.URL.Path {
			r.handlers[i].handler(response, request)
			return
		}
	}
	http.NotFound(response, request)
}

func (r *SliceRouter) AddHandler(path string, handler HandlerFunction) {

	if r.handlers == nil {
		r.handlers = make([]*routeHandler, 0, 10)
	}

	rh := &routeHandler{path: path, handler: handler}
	r.handlers = append(r.handlers, rh)
}

type MapRouter struct {
	handlers map[string]HandlerFunction
}

func (r *MapRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, found := r.handlers[request.URL.Path]
	if !found {
		http.NotFound(response, request)
		return
	}

	handler(response, request)
}

func (r *MapRouter) AddHandler(path string, handler HandlerFunction) {

	if r.handlers == nil {
		r.handlers = make(map[string]HandlerFunction)
	}

	r.handlers[path] = handler
}

type PrefixTreeRouter struct {
	tree Tree
}

func (r *PrefixTreeRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler := r.tree.Find(request.URL.Path)
	if handler == nil {
		http.NotFound(response, request)
		return
	}

	handler(response, request)
}

func (r *PrefixTreeRouter) AddHandler(path string, handler HandlerFunction) {
	r.tree.Insert(path, handler)
}
