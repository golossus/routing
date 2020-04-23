package routing

import (
	"context"
	"fmt"
	"net/http"
	"strings"
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

	leaf := ctx.(*node)

	path := request.URL.Path
	return buildURLParameters(leaf, path, len(path), 0)
}

func buildURLParameters(leaf *node, path string, offset int, paramsCount uint) URLParameterBag {

	if leaf == nil {
		return newURLParameterBag(paramsCount, false)
	}

	var paramsBag URLParameterBag

	if leaf.t == nodeTypeDynamic {
		start := strings.LastIndex(path[:offset], leaf.parent.prefix) + len(leaf.parent.prefix)
		paramsBag = buildURLParameters(leaf.parent, path, start, paramsCount+1)
		paramsBag.add(leaf.prefix, path[start:offset])
	} else {
		paramsBag = buildURLParameters(leaf.parent, path, offset-len(leaf.prefix), paramsCount)
	}

	return paramsBag
}

// Router is a structure where all routes are stored
type Router struct {
	trees  map[string]*tree
	asName string
	routes map[string]*node
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

	leaf := tree.find(request.URL.Path)
	if leaf == nil {
		http.NotFound(response, request)
		return
	}

	ctx := context.Background()
	leaf.handler(response, request.WithContext(context.WithValue(ctx, ctxKey, leaf)))
}

// As method sets a name for the next registered route
func (r *Router) As(asName string)  *Router {
	r.asName = asName
	return r
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

	if nil == r.routes {
		r.routes = make(map[string]*node)
	}

	if _, ok := r.trees[verb]; !ok {
		r.trees[verb] = &tree{}
	}

	leaf := r.trees[verb].insert(parser.chunks, handler)

	if r.asName != "" {
		r.routes[r.asName] = leaf
		r.asName = ""
	}

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
		t.root.parent = leafNew

		if leafNew.t == nodeTypeDynamic {
			leafNew.stops[t.root.prefix[0]] = t.root
		} else {
			leafNew.child = t.root
		}

		r.trees[verb].root = combine(r.trees[verb].root, rootNew)
	}

	for name, leaf := range router.routes {
		r.routes[fmt.Sprintf("%s%s", r.asName, name)] = leaf
	}
	r.asName = ""

	return nil
}

// GenerateURL generates a URL from route name
func (r *Router) GenerateURL(name string, params URLParameterBag) (string, error) {
	node, ok := r.routes[name]
	if !ok {
		return "", fmt.Errorf("route name %s not found", name)
	}

	url := ""
	for node != nil {

		if node.t == nodeTypeStatic {
			url = node.prefix + url
		}else{
			p, err := params.GetByName(node.prefix)
			if err != nil {
				return p, err
			}
			if node.regexp != nil && !node.regexp.MatchString(p) {
				return "", fmt.Errorf("param %s with value %s is not valid", node.prefix, p)
			}
			url = p + url
		}

		node = node.parent
	}
	return url, nil
}
