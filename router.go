package routing

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type paramsKey int

var ctxKey paramsKey

var handlers = make(map[string]http.HandlerFunc)

// AddHandler adds an http.HandlerFunc into a list of handlers to be retrieved
// by name (canonical or alias) on runtime
func AddHandler(handler http.HandlerFunc, aliases ...string) {
	name := runtime.FuncForPC(reflect.ValueOf(handler).Pointer()).Name()
	handlers[strings.TrimRight(name, "-fm")] = handler

	for _, alias := range aliases {
		handlers[alias] = handler
	}
}

// GetHandler retrieves an http.HandlerFunc given a name from the list of handlers
func GetHandler(name string) (http.HandlerFunc, error) {
	handler, ok := handlers[name]
	if !ok {
		return nil, fmt.Errorf("handler with name %s not registered", name)
	}

	return handler, nil
}

// GetURLParameters is in charge of retrieve dynamic parameter of the URL within your route.
// For example, User's ID in /users/{userId}
func GetURLParameters(request *http.Request) URLParameterBag {
	ctx := request.Context().Value(ctxKey)
	if ctx == nil {
		return newURLParameterBag(0)
	}

	leaf := ctx.(nodeInterface)

	path := request.URL.Path
	return buildURLParameters(leaf, path, len(path), 0)
}

func buildURLParameters(leaf nodeInterface, path string, offset int, paramsCount uint) URLParameterBag {

	if leaf == nil {
		return newURLParameterBag(paramsCount)
	}

	var paramsBag URLParameterBag

	switch leaf.(type) {
	case *nodeDynamic:
		start := strings.LastIndex(path[:offset], leaf.getParent().getPrefix()) + len(leaf.getParent().getPrefix())
		paramsBag = buildURLParameters(leaf.getParent(), path, start, paramsCount+1)
		paramsBag.add(leaf.getPrefix(), path[start:offset])
	case *nodeStatic:
		paramsBag = buildURLParameters(leaf.getParent(), path, offset-len(leaf.getPrefix()), paramsCount)
	}

	return paramsBag
}

// Router is a structure where all routes are stored
type Router struct {
	trees  map[string]*tree
	asName string
	routes map[string]nodeInterface
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

	if leaf.hasParameters() {
		ctx := context.Background()
		request = request.WithContext(context.WithValue(ctx, ctxKey, leaf))
	}
	leaf.handler()(response, request)
}

// As method sets a name for the next registered route
func (r *Router) As(asName string) *Router {
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
		r.routes = make(map[string]nodeInterface)
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
		t.root.setParent(leafNew)

		switch leafNew.(type) {
		case *nodeDynamic:
			leafNew.(*nodeDynamic).childrenNodes[t.root.getPrefix()[0]] = t.root
		case *nodeStatic:
			leafNew.(*nodeStatic).childNode = t.root
		}

		if r.trees[verb].root == nil {
			r.trees[verb].root = rootNew
		}else{
			r.trees[verb].root = r.trees[verb].root.merge(rootNew)
		}
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

	var url strings.Builder
	err := getUri(node, &url, params)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func getUri(node nodeInterface, url *strings.Builder, params URLParameterBag) error {

	if node == nil{
		return nil
	}

	err := getUri(node.getParent(), url, params)
	if err != nil{
		return err
	}

	switch node.(type) {
	case *nodeStatic:
		url.WriteString(node.getPrefix())
	case *nodeDynamic:
		p, err := params.GetByName(node.getPrefix())
		if err != nil {
			return err
		}
		if node.(*nodeDynamic).regexp != nil && !node.(*nodeDynamic).regexp.MatchString(p) {
			return fmt.Errorf("param %s with value %s is not valid", node.getPrefix(), p)
		}
		url.WriteString(p)
	}

	return nil
}

// RouteDef defines a route definition
type RouteDef struct {
	Method  string
	Schema  string
	Handler string
	Name    string
}

// Loader loads a list routes
type Loader interface {
	Load() []RouteDef
}

// Load registers a list of routes retrieved from a loader
func (r *Router) Load(loader Loader) error {
	for _, route := range loader.Load() {
		handler, err := GetHandler(route.Handler)
		if err != nil {
			return err
		}
		if route.Name != "" {
			r.As(route.Name)
		}
		err = r.Register(route.Method, route.Schema, handler)
		if err != nil {
			return err
		}
	}
	return nil
}

// Load registers a list of routes retrieved from a loader
func (r *Router) PrioritizeByWeight() {
	for _, tree := range r.trees {
		_ = calcWeight(tree.root)
		tree.root = sortByWeight(tree.root)
	}
}
