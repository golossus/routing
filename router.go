package http_router

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

type Router interface {
	http.Handler
	AddHandler(string, string, HandlerFunction)
}

type TreeRouter struct {
	tree Tree
}

func (r *TreeRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	handler, params := r.tree.Find(request.Method, request.URL.Path)
	if handler == nil {
		http.NotFound(response, request)
		return
	}

	ctx := context.Background()
	handler(response, request.WithContext(context.WithValue(ctx, ctxKey, params)))
}

func (r *TreeRouter) Head(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodHead, path, handler)
}

func (r *TreeRouter) Get(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodGet, path, handler)
}

func (r *TreeRouter) Post(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPost, path, handler)
}

func (r *TreeRouter) Put(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPut, path, handler)
}

func (r *TreeRouter) Patch(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodPatch, path, handler)
}

func (r *TreeRouter) Delete(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodDelete, path, handler)
}

func (r *TreeRouter) Connect(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodConnect, path, handler)
}

func (r *TreeRouter) Options(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodOptions, path, handler)
}

func (r *TreeRouter) Trace(path string, handler HandlerFunction) {
	r.AddHandler(http.MethodTrace, path, handler)
}

func (r *TreeRouter) Any(path string, handler HandlerFunction) {
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

func (r *TreeRouter) AddHandler(verb, path string, handler HandlerFunction) {
	parser := newParser(path)
	_, err := parser.parse()
	if err != nil {
		panic(err)
	}

	r.tree.Insert(verb, parser.chunks, handler)
}
