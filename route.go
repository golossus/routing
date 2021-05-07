package routing

import "net/http"

// route represents a single route configuration
type route struct {
	method  string
	path    string
	handler http.HandlerFunc
	options MatchingOptions
}

// routeBuilder defines an interface to build a route for a Router
type routeBuilder struct {
	r    *Router
	curr route
}

// Register registers the current route in the corresponding Router
func (b *routeBuilder) Register() error {
	return b.r.Register(b.curr.method, b.curr.path, b.curr.handler, b.curr.options)
}

// Method sets the http method or verb of the current route
func (b *routeBuilder) Method(method string) *routeBuilder {
	b.curr.method = method
	return b
}

// Path sets the path schema to match of the current route
func (b *routeBuilder) Path(path string) *routeBuilder {
	b.curr.path = path
	return b
}

// Handler sets the handler to execute of the current route when matched
func (b *routeBuilder) Handler(handler http.HandlerFunc) *routeBuilder {
	b.curr.handler = handler
	return b
}

// Name sets the name of the current route to be used to generate the Url
func (b *routeBuilder) Name(name string) *routeBuilder {
	b.curr.options.Name = name
	return b
}

// Host sets the host name schema to match of the current route
func (b *routeBuilder) Host(host string) *routeBuilder {
	b.curr.options.Host = host
	return b
}

// Schemas set the list of Url schemas to match against
func (b *routeBuilder) Schemas(schemas ...string) *routeBuilder {
	b.curr.options.Schemas = schemas
	return b
}

// Header adds in the current route a header name and value to match against
func (b *routeBuilder) Header(header, value string) *routeBuilder {
	if b.curr.options.Headers == nil {
		b.curr.options.Headers = make(map[string]string)
	}

	b.curr.options.Headers[header] = value
	return b
}

// QueryParam adds in the current route a query string parameter name and value
// to match against
func (b *routeBuilder) QueryParam(param, value string) *routeBuilder {
	if b.curr.options.QueryParams == nil {
		b.curr.options.QueryParams = make(map[string]string)
	}

	b.curr.options.QueryParams[param] = value
	return b
}

// Matcher sets a CustomMatcher function to match the route
func (b *routeBuilder) Matcher(f CustomMatcher) *routeBuilder {
	b.curr.options.Custom = f
	return b
}
