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

	leaf := ctx.(*node)

	path := request.URL.Path
	return buildURLParameters(leaf, path, len(path), 0)
}

func buildURLParameters(leaf *node, path string, offset int, paramsCount uint) URLParameterBag {

	if leaf == nil {
		return newURLParameterBag(paramsCount)
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

	leaf := tree.find(request)
	if leaf == nil {
		http.NotFound(response, request)
		return
	}

	if leaf.hasParameters() {
		ctx := context.Background()
		request = request.WithContext(context.WithValue(ctx, ctxKey, leaf))
	}
	leaf.handler(response, request)
}

// As method sets a name for the next registered route.
//
// Deprecated: MatchingOptions should be used instead and will have preference
// over this method. It will be deleted on minor version after v1.2.0
func (r *Router) As(asName string) *Router {
	r.asName = asName
	return r
}

type MatchingOptions struct {
	Name string
	Host string
}

func NewMatchingOptions() MatchingOptions {
	return MatchingOptions{
		Name: "",
		Host: "",
	}
}

// Register adds a new route in the router
func (r *Router) Register(verb, path string, handler http.HandlerFunc, options ...MatchingOptions) error {
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


	rname := r.asName
	if len(options) > 0 {
		rname = options[0].Name

		if options[0].Host != "" {
			leaf.matchers = append(leaf.matchers, byHost(options[0].Host))
		}
	}

	if rname != "" {
		r.routes[rname] = leaf
	}
	r.asName = ""

	return nil
}

// Head is a method to register a new HEAD route in the router.
func (r *Router) Head(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodHead, path, handler, options...)
}

// Get is a method to register a new GET route in the router.
func (r *Router) Get(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodGet, path, handler, options...)
}

// Post is a method to register a new POST route in the router.
func (r *Router) Post(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodPost, path, handler, options...)
}

// Put is a method to register a new PUT route in the router.
func (r *Router) Put(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodPut, path, handler, options...)
}

// Patch is a method to register a new PATCH route in the router.
func (r *Router) Patch(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodPatch, path, handler, options...)
}

// Delete is a method to register a new DELETE route in the router.
func (r *Router) Delete(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodDelete, path, handler, options...)
}

// Connect is a method to register a new CONNECT route in the router.
func (r *Router) Connect(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodConnect, path, handler, options...)
}

// Options is a method to register a new OPTIONS route in the router.
func (r *Router) Options(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodOptions, path, handler, options...)
}

// Trace is a method to register a new TRACE route in the router.
func (r *Router) Trace(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	return r.Register(http.MethodTrace, path, handler, options...)
}

// Any is a method to register a new route with all the verbs.
func (r *Router) Any(path string, handler http.HandlerFunc, options ...MatchingOptions) error {
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
		if err := r.Register(verb, path, handler, options...); err != nil {
			return err
		}
	}

	return nil
}

// Prefix combines two routers under a custom path prefix
func (r *Router) Prefix(path string, router *Router, name ...string) error {
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

	rname := r.asName
	if len(name) > 0 {
		rname = name[0]
	}

	for name, leaf := range router.routes {
		r.routes[fmt.Sprintf("%s%s", rname, name)] = leaf
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

func getUri(node *node, url *strings.Builder, params URLParameterBag) error {

	if node == nil {
		return nil
	}

	err := getUri(node.parent, url, params)
	if err != nil {
		return err
	}

	if node.t == nodeTypeStatic {
		url.WriteString(node.prefix)
	} else {
		p, err := params.GetByName(node.prefix)
		if err != nil {
			return err
		}
		if node.regexp != nil && !node.regexp.MatchString(p) {
			return fmt.Errorf("param %s with value %s is not valid", node.prefix, p)
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
