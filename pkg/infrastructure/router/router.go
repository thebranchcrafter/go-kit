package router

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Middleware func(gin.HandlerFunc) gin.HandlerFunc

// Route represents an HTTP route with its method, path, handler, and middleware.
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []Middleware
}

type Router interface {
	Handle(method, path string, handler http.HandlerFunc, middleware ...Middleware)
	Serve(addr string) error
	Handler() http.Handler
	ServeStatic(url, absPath string) http.Handler
}
