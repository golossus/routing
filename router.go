package routing

import (
	"context"
	"net/http"
)

type paramsKey int

var ctxKey paramsKey

// GetURLParameters is in charge of retrieve dynamic parameter of the URL within your route. For example, User's ID in /users/{userId}
func GetURLParameters(request *http.Request) URLParameterBag {
	return request.Context().Value(ctxKey).(URLParameterBag)
}

// HandlerFunction is a function to be executed when the router match a route
type HandlerFunction func(http.ResponseWriter, *http.Request)

// Router is a structure where all routes are stored
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

// Head is a method to register a new HEAD route in the router.
func (r *Router) Head(path string, handler HandlerFunction) {
	r.addHandler(http.MethodHead, path, handler)
}

// Get is a method to register a new GET route in the router.
func (r *Router) Get(path string, handler HandlerFunction) {
	r.addHandler(http.MethodGet, path, handler)
}

// Post is a method to register a new POST route in the router.
func (r *Router) Post(path string, handler HandlerFunction) {
	r.addHandler(http.MethodPost, path, handler)
}

// Put is a method to register a new PUT route in the router.
func (r *Router) Put(path string, handler HandlerFunction) {
	r.addHandler(http.MethodPut, path, handler)
}

// Patch is a method to register a new PATCH route in the router.
func (r *Router) Patch(path string, handler HandlerFunction) {
	r.addHandler(http.MethodPatch, path, handler)
}

// Delete is a method to register a new DELETE route in the router.
func (r *Router) Delete(path string, handler HandlerFunction) {
	r.addHandler(http.MethodDelete, path, handler)
}

// Connect is a method to register a new CONNECT route in the router.
func (r *Router) Connect(path string, handler HandlerFunction) {
	r.addHandler(http.MethodConnect, path, handler)
}

// Options is a method to register a new OPTIONS route in the router.
func (r *Router) Options(path string, handler HandlerFunction) {
	r.addHandler(http.MethodOptions, path, handler)
}

// Trace is a method to register a new TRACE route in the router.
func (r *Router) Trace(path string, handler HandlerFunction) {
	r.addHandler(http.MethodTrace, path, handler)
}

// Any is a method to register a new route with all the verbs.
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
		r.addHandler(verb, path, handler)
	}
}

func (r *Router) addHandler(verb, path string, handler HandlerFunction) {
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
