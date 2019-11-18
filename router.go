package hw14_go

import (
	"net/http"
)

type HandlerFunction func(http.ResponseWriter, *http.Request)

type routeHandler struct {
	path    string
	handler HandlerFunction
}

type Router struct {
	handlers []*routeHandler
}

func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	for i := 0; i < len(r.handlers); i++ {
		if r.handlers[i].path == request.URL.Path {
			r.handlers[i].handler(response, request)
		}
	}
	http.NotFound(response, request)
}

func (r *Router) AddHandler(path string, handler HandlerFunction) {

	if r.handlers == nil {
		r.handlers = make([]*routeHandler, 0, 10)
	}

	rh := &routeHandler{path: path, handler: handler}
	r.handlers = append(r.handlers, rh)
}
