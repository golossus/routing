package http_router

import (
	"context"
	"net/http"
)

type HandlerFunction func(http.ResponseWriter, *http.Request)

type Router interface {
	http.Handler
	AddHandler(string, string, HandlerFunction)
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

func (r *SliceRouter) AddHandler(verb string, path string, handler HandlerFunction) {

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

func (r *MapRouter) AddHandler(verb string, path string, handler HandlerFunction) {

	if r.handlers == nil {
		r.handlers = make(map[string]HandlerFunction)
	}

	r.handlers[path] = handler
}

type PrefixTreeRouter struct {
	tree Tree
}

func (r *PrefixTreeRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params := r.tree.Find(request.Method, request.URL.Path)
	if handler == nil {
		http.NotFound(response, request)
		return
	}

	ctx := context.Background()
	handler(response, request.WithContext(context.WithValue(ctx, ParamsBagKey, params)))
}

func (r *PrefixTreeRouter) Head(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodHead, path, handler)
}

func (r *PrefixTreeRouter) Get(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodGet, path, handler)
}

func (r *PrefixTreeRouter) Post(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPost, path, handler)
}

func (r *PrefixTreeRouter) Put(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPut, path, handler)
}

func (r *PrefixTreeRouter) Patch(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPatch, path, handler)
}

func (r *PrefixTreeRouter) Delete(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodDelete, path, handler)
}

func (r *PrefixTreeRouter) Connect(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodConnect, path, handler)
}

func (r *PrefixTreeRouter) Options(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodOptions, path, handler)
}

func (r *PrefixTreeRouter) Trace(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodTrace, path, handler)
}

func (r *PrefixTreeRouter) Any(path string, handler HandlerFunction) {
	kvs := [9]string{
		http.MethodHead,
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
	for _, verb := range kvs {
		r.AddHandler(verb, path, handler)
	}
}

func (r *PrefixTreeRouter) AddHandler(verb, path string, handler HandlerFunction) {
	parser := NewParser(path)
	_, err := parser.parse()
	if err != nil {
		panic(err)
	}

	r.tree.Insert(verb, parser.chunks, handler)
}