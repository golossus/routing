package routing

import (
	"context"
	"net/http"
)

type paramsKey int

var ctxKey paramsKey

// GetURLParameters is in charge of retrieve dynamic parameter of the URL within your route.
// For example, User's ID in /users/{userId}
func GetURLParameters(request *http.Request) URLParameterBag {
	ctx := request.Context().Value(ctxKey)
	if ctx == nil {
		return newURLParameterBag(0, true)
	}
	return ctx.(URLParameterBag)
}

// Router is a structure where all routes are stored
type Router struct {
	trees map[string]*tree
}

// NewRouter returns an empty Router
func NewRouter() Router {
	return Router{}
}

// ServerHTTP executes the HandlerFunc if the request path is found
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
func (r *Router) Head(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodHead, path, handler)
}

// Get is a method to register a new GET route in the router.
func (r *Router) Get(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodGet, path, handler)
}

// Post is a method to register a new POST route in the router.
func (r *Router) Post(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodPost, path, handler)
}

// Put is a method to register a new PUT route in the router.
func (r *Router) Put(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodPut, path, handler)
}

// Patch is a method to register a new PATCH route in the router.
func (r *Router) Patch(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodPatch, path, handler)
}

// Delete is a method to register a new DELETE route in the router.
func (r *Router) Delete(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodDelete, path, handler)
}

// Connect is a method to register a new CONNECT route in the router.
func (r *Router) Connect(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodConnect, path, handler)
}

// Options is a method to register a new OPTIONS route in the router.
func (r *Router) Options(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodOptions, path, handler)
}

// Trace is a method to register a new TRACE route in the router.
func (r *Router) Trace(path string, handler http.HandlerFunc) error {
	return r.Register(http.MethodTrace, path, handler)
}

// Any is a method to register a new route with all the verbs.
func (r *Router) Any(path string, handler http.HandlerFunc) error {
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
		if err := r.Register(verb, path, handler); err != nil {
			return err
		}
	}

	return nil
}

// Register adds a new route in the router
func (r *Router) Register(verb, path string, handler http.HandlerFunc) error {
	parser := newParser(path)
	_, err := parser.parse()
	if err != nil {
		return err
	}

	if nil == r.trees {
		r.trees = make(map[string]*tree)
	}

	if _, ok := r.trees[verb]; !ok {
		r.trees[verb] = &tree{}
	}

	r.trees[verb].insert(parser.chunks, handler)

	return nil
}

// Prefix combines two routers under a custom path prefix
func (r *Router) Prefix(path string, router *Router) error {
	parser := newParser(path)
	_, err := parser.parse()
	if err != nil {
		return err
	}

	if nil == r.trees {
		r.trees = make(map[string]*tree)
	}

	for verb, t := range router.trees {
		if _, ok := r.trees[verb]; !ok {
			r.trees[verb] = &tree{}
		}

		rootNew, leafNew := createTreeFromChunks(parser.chunks)

		if leafNew.t == nodeTypeDynamic {
			leafNew.stops[t.root.prefix[0]] = t.root
		} else {
			leafNew.child = t.root
		}

		r.trees[verb].root = combine(r.trees[verb].root, rootNew)
	}

	return nil
}
