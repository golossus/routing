package routing

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type MiddlewarePipe struct {
	middlewares []Middleware
}

func NewMiddlewarePipe() *MiddlewarePipe {
	return &MiddlewarePipe{}
}

func (s *MiddlewarePipe) Next(middleware ...Middleware) *MiddlewarePipe {
	s.middlewares = append(s.middlewares, middleware...)
	return s
}

func (s *MiddlewarePipe) Pipe(pipe *MiddlewarePipe) {
	s.middlewares = append(s.middlewares, pipe.middlewares...)
}

func (s *MiddlewarePipe) Then(next http.HandlerFunc) http.HandlerFunc {
	for j := len(s.middlewares)-1; j >= 0; j-- {
		next = s.middlewares[j](next)
	}

	return next
}
