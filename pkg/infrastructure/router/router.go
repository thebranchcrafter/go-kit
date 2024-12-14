package router

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Route represents an HTTP route with its method, path, handler, and middleware.
type Route struct {
	Method      string
	Path        string
	Handler     http.HandlerFunc
	Middlewares []Middleware
}

type Router interface {
	Handle(method, path string, handler http.HandlerFunc, middleware ...Middleware)
	Serve(addr string) error
	Handler() http.Handler
}
