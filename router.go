package routing

import (
	"context"
	"net/http"
)

type paramsKey int

var ctxKey paramsKey

func GetUrlParameters(request *http.Request) UrlParameterBag {
	return request.Context().Value(ctxKey).(UrlParameterBag)
}

type HandlerFunction func(http.ResponseWriter, *http.Request)

type Router struct {
	trees map[string]*tree
}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
  tree, ok := r.trees[request.Method]
	if !ok {
		http.NotFound(response, request)
		return
	}
  
  handler, params := tree.find(request.URL.Path)
	if handler == nil {
		http.NotFound(response, request)
		return
	}

	ctx := context.Background()
	handler(response, request.WithContext(context.WithValue(ctx, ctxKey, params)))
}

func (r *Router) Head(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodHead, path, handler)
}

func (r *Router) Get(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodGet, path, handler)
}

func (r *Router) Post(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPost, path, handler)
}

func (r *Router) Put(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPut, path, handler)
}

func (r *Router) Patch(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPatch, path, handler)
}

func (r *Router) Delete(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodDelete, path, handler)
}

func (r *Router) Connect(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodConnect, path, handler)
}

func (r *Router) Options(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodOptions, path, handler)
}

func (r *Router) Trace(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodTrace, path, handler)
}

func (r *Router) Any(path string, handler HandlerFunction) {
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

func (r *Router) AddHandler(verb, path string, handler HandlerFunction) {
	parser := newParser(path)
	_, err := parser.parse()
	if err != nil {
		panic(err)
	}

	if nil == r.trees {
		r.trees = make(map[string]*tree)
	}

	if _, ok := r.trees[verb]; !ok {
		r.trees[verb] = &tree{}
	}

	r.trees[verb].insert(parser.chunks, handler)
}
