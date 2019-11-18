package hw14_go

import (
	"fmt"
	"net/http"
)

func Greet() {
	fmt.Println("Hola HW14")
}

type HandlerFunction func(http.ResponseWriter, *http.Request)


type RouteHandler struct {
	path string
	handler HandlerFunction
}

type Router struct {
	routeHandlers []*RouteHandler

}

func (r *Router) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("implement me")
}

func (r *Router) AddHandler(path string, handler HandlerFunction) {

	if (r.routeHandlers == nil) {
		r.routeHandlers = make([]*RouteHandler, 0, 10)
	}

	routeHandler := &RouteHandler { path: path, handler: handler }
	r.routeHandlers = append(r.routeHandlers, routeHandler)
}
