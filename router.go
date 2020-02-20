package http_router

import (
	"context"
	"net/http"
)

const (
	GET     = "GET"
	HEAD    = "HEAD"
	POST    = "POST"
	PUT     = "PUT"
	PATCH   = "PATCH"
	DELETE  = "DELETE"
	CONNECT = "CONNECT"
	OPTIONS = "OPTIONS"
	TRACE   = "TRACE"
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

const (
	ParamsBagKey = "urlParameters"
)

func (r *PrefixTreeRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params := r.tree.Find(request.Method, request.URL.Path)
	if handler == nil {
		http.NotFound(response, request)
		return
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, ParamsBagKey, params)
	handler(response, request.WithContext(ctx))
}

func (r *PrefixTreeRouter) Head(path string, handler HandlerFunction) {
	r.AddHandler(HEAD, path, handler)
}

func (r *PrefixTreeRouter) Get(path string, handler HandlerFunction) {
	r.AddHandler(GET, path, handler)
}

func (r *PrefixTreeRouter) Post(path string, handler HandlerFunction) {
	r.AddHandler(POST, path, handler)
}

func (r *PrefixTreeRouter) Put(path string, handler HandlerFunction) {
	r.AddHandler(PUT, path, handler)
}

func (r *PrefixTreeRouter) Patch(path string, handler HandlerFunction) {
	r.AddHandler(PATCH, path, handler)
}

func (r *PrefixTreeRouter) Delete(path string, handler HandlerFunction) {
	r.AddHandler(DELETE, path, handler)
}

func (r *PrefixTreeRouter) Connect(path string, handler HandlerFunction) {
	r.AddHandler(CONNECT, path, handler)
}

func (r *PrefixTreeRouter) Options(path string, handler HandlerFunction) {
	r.AddHandler(OPTIONS, path, handler)
}

func (r *PrefixTreeRouter) Trace(path string, handler HandlerFunction) {
	r.AddHandler(TRACE, path, handler)
}

func (r *PrefixTreeRouter) Any(path string, handler HandlerFunction) {
	kvs := map[string]string{HEAD: HEAD, GET: GET, POST: POST, PUT: PUT, PATCH: PATCH, DELETE: DELETE, CONNECT: CONNECT, OPTIONS: OPTIONS, TRACE: TRACE}
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

type urlParameter struct {
	name  string
	value string
}

type UrlParameterBag struct {
	params []urlParameter
}

func (u *UrlParameterBag) addParameter(param urlParameter) {
	if u.params == nil {
		u.params = make([]urlParameter, 0, 5)
	}

	u.params = append(u.params, param)
}

func (u *UrlParameterBag) GetByName(name string, def string) string {
	for _, item := range u.params {
		if item.name == name {
			return item.value
		}
	}

	return def
}

func (u *UrlParameterBag) GetByIndex(index uint, def string) string {
	i := int(index)
	if len(u.params) <= i {
		return def
	}

	return u.params[i].value
}

func NewUrlParameterBag() UrlParameterBag {
	return UrlParameterBag{}
}
