package routing

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
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

var matchers = make(map[string]CustomMatcher)

// AddCustomMatcher adds a customer route matcher into a list of matchers to be
// retrieved by name (canonical or alias) on runtime
func AddCustomMatcher(m CustomMatcher, aliases ...string) {
	name := runtime.FuncForPC(reflect.ValueOf(m).Pointer()).Name()
	matchers[strings.TrimRight(name, "-fm")] = m

	for _, alias := range aliases {
		matchers[alias] = m
	}
}

// GetCustomMatcher retrieves a custom matcher given a name from the list of matchers
func GetCustomMatcher(name string) (CustomMatcher, error) {
	m, ok := matchers[name]
	if !ok {
		return nil, fmt.Errorf("custom matcher with name %s not registered", name)
	}

	return m, nil
}

// GetURLParameters is in charge of retrieve dynamic parameter of the URL within your route.
// For example, User's ID in /users/{userId}
func GetURLParameters(request *http.Request) URLParameterBag {
	ctx := request.Context().Value(ctxKey)
	if ctx == nil {
		return newURLParameterBag(0)
	}

	leaf := ctx.(*node)

	urlParams := buildURLParameters(leaf, request.URL.Path, len(request.URL.Path), 0)

	for _, matcher := range leaf.matchers {
		if matches, hostLeaf := matcher(request); matches {
			urlParams = urlParams.merge(buildURLParameters(hostLeaf, request.Host, len(request.Host), 0))
		}
	}

	return urlParams
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

// RouterConfig is a structure to set the router configuration
type RouterConfig struct {
	EnableAutoMethodHead           bool
	EnableAutoMethodOptions        bool
	EnableMethodNotAllowedResponse bool
}

// Router is a structure where all routes are stored
type Router struct {
	config RouterConfig
	trees  map[string]*tree
	asName string
	routes map[string]*node
}

// NewRouter returns an empty Router
func NewRouter(configs ...RouterConfig) Router {
	var defaultConfig RouterConfig
	if len(configs) > 0 {
		defaultConfig = configs[0]
	}
	return Router{config: defaultConfig}
}

// ServerHTTP executes the HandlerFunc if the request path is found
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	tree, ok := r.trees[request.Method]
	if !ok {
		r.notFoundOrMethodNotAllowed(response, request)
		return
	}

	leaf := tree.find(request)
	if leaf == nil {
		r.notFoundOrMethodNotAllowed(response, request)
		return
	}

	if leaf.hasParameters() {
		ctx := context.Background()
		request = request.WithContext(context.WithValue(ctx, ctxKey, leaf))
	}
	leaf.handler(response, request)
}

func (r *Router) notFoundOrMethodNotAllowed(response http.ResponseWriter, request *http.Request) {
	if !r.config.EnableMethodNotAllowedResponse {
		http.NotFound(response, request)
		return
	}

	availVerbs := getAvailableMethods(r, request)
	if len(availVerbs) == 0 {
		http.NotFound(response, request)
		return
	}

	response.Header().Set("Allow", strings.Join(availVerbs, ", "))
	http.Error(response, "405 method not allowed", http.StatusMethodNotAllowed)
}

// As method sets a name for the next registered route.
//
// Deprecated: MatchingOptions should be used instead and will have preference
// over this method. It will be deleted on minor version after v1.2.0
func (r *Router) As(asName string) *Router {
	r.asName = asName
	return r
}

// MatchingOptions is a structure to define a route name and extend the matching options
type MatchingOptions struct {
	Name        string
	Host        string
	Schemas     []string
	Headers     map[string]string
	QueryParams map[string]string
	Custom      CustomMatcher
}

// NewMatchingOptions returns the MatchingOptions structure
func NewMatchingOptions() MatchingOptions {
	return MatchingOptions{
		Name:        "",
		Host:        "",
		Schemas:     nil,
		Headers:     map[string]string{},
		QueryParams: map[string]string{},
		Custom:      nil,
	}
}

// Register adds a new route in the router
func (r *Router) Register(verb, path string, handler http.HandlerFunc, options ...MatchingOptions) error {
	if len(verb) < 3 {
		return fmt.Errorf("invalid verb %s", verb)
	}

	if handler == nil {
		return fmt.Errorf("handler can not be nil")
	}

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
	r.asName = ""

	if len(options) > 0 {
		rname = options[0].Name

		if options[0].Host != "" {
			matcherByHost, err := byHost(options[0].Host)
			if err != nil {
				return err
			}
			leaf.matchers = append(leaf.matchers, matcherByHost)
		}

		if len(options[0].Schemas) > 0 {
			matcherBySchemas, err := bySchemas(options[0].Schemas...)
			if err != nil {
				return err
			}
			leaf.matchers = append(leaf.matchers, matcherBySchemas)
		}

		if len(options[0].Headers) > 0 {
			matcherByHeaders := byHeaders(options[0].Headers)
			leaf.matchers = append(leaf.matchers, matcherByHeaders)
		}

		if len(options[0].QueryParams) > 0 {
			matcherByQueryParams := byQueryParameters(options[0].QueryParams)
			leaf.matchers = append(leaf.matchers, matcherByQueryParams)
		}

		if options[0].Custom != nil {
			matcherByCustomFunc := byCustomMatcher(options[0].Custom)
			leaf.matchers = append(leaf.matchers, matcherByCustomFunc)
		}
	}

	rname = r.generateRouteName(rname, parser)

	r.routes[rname] = leaf

	if r.config.EnableAutoMethodHead && verb == http.MethodGet {
		_ = r.Register(http.MethodHead, path, handler, options...)
	}

	if r.config.EnableAutoMethodOptions && verb != http.MethodOptions {
		_ = r.Register(http.MethodOptions, path, getAutoMethodOptionsHandler(r), options...)
	}

	return nil
}

func (r *Router) generateRouteName(baseName string, parser *parser) string {
	if baseName == "" && parser != nil {
		for _, c := range parser.chunks {
			baseName += c.v
		}
		baseName = strings.Trim(strings.Replace(baseName, "/", "_", -1), "_")
	}

	i := 0
	existsName := baseName
	for {
		_, ok := r.routes[existsName]
		if !ok {
			break
		}
		i++
		existsName = baseName + "_" + strconv.Itoa(i)
	}

	return existsName
}

// NewRoute is a method to register a route in the router through a builder interface.
func (r *Router) NewRoute() *routeBuilder {
	return &routeBuilder{r, route{options: MatchingOptions{}}}
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

	r.asName = ""

	for name, leaf := range router.routes {
		r.routes[r.generateRouteName(name, nil)] = leaf
	}

	return nil
}

// StaticFiles  will serve files from a directory under a prefix path
func (r *Router) StaticFiles(prefix, dir string) error {
	return r.Register("GET", prefix+"/{name:.*}", func(writer http.ResponseWriter, request *http.Request) {

		urlParams := GetURLParameters(request)
		name, _ := urlParams.GetByName("name")

		request.URL.Path = name
		http.FileServer(http.Dir(dir)).ServeHTTP(writer, request)
	})
}

// Redirect will redirect a path to an url
func (r *Router) Redirect(path, url string, code ...int) error {
	return r.Register(http.MethodGet, path, getRedirectHandler(url, code...))
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

// RouteDef represents a route definition
type RouteDef struct {
	Method  string
	Path    string
	Handler string
	Options RouteDefOptions
}

// RouteDefOptions represents a route definition extra options
type RouteDefOptions struct {
	Name          string
	Host          string
	Schemas       []string
	Headers       map[string]string
	QueryParams   map[string]string
	CustomMatcher string
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

		var matcher CustomMatcher
		if len(route.Options.CustomMatcher) > 0 {
			matcher, err = GetCustomMatcher(route.Options.CustomMatcher)
			if err != nil {
				return err
			}
		}

		options := MatchingOptions{
			Name:        route.Options.Name,
			Host:        route.Options.Host,
			Schemas:     route.Options.Schemas,
			Headers:     route.Options.Headers,
			QueryParams: route.Options.QueryParams,
			Custom:      matcher,
		}
		err = r.Register(route.Method, route.Path, handler, options)
		if err != nil {
			return err
		}
	}
	return nil
}

// PrioritizeByWeight changes the router underlying tree to prioritize search
// through the branches of higher weight.
func (r *Router) PrioritizeByWeight() {
	for _, tree := range r.trees {
		_ = calcWeight(tree.root)
		tree.root = sortByWeight(tree.root)
	}
}
