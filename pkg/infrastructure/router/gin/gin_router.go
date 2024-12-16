package gin_router

import (
	"github.com/gin-gonic/gin"
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/router"
	"net/http"
)

type GinRouter struct {
	engine *gin.Engine
}

func NewGinRouter() *GinRouter {
	engine := gin.Default()
	return &GinRouter{engine: engine}
}

func (g *GinRouter) Handle(method, path string, handler gin.HandlerFunc, middleware ...router.Middleware) {
	// Wrap the gin.HandlerFunc with middleware
	finalHandler := handler
	for _, m := range middleware {
		finalHandler = m(finalHandler)
	}

	// Register the route
	switch method {
	case http.MethodGet:
		g.engine.GET(path, finalHandler)
	case http.MethodPost:
		g.engine.POST(path, finalHandler)
	case http.MethodPut:
		g.engine.PUT(path, finalHandler)
	case http.MethodDelete:
		g.engine.DELETE(path, finalHandler)
	default:
		panic("unsupported HTTP method: " + method)
	}
}

func (g *GinRouter) Serve(addr string) error {
	return g.engine.Run(addr)
}

func (g *GinRouter) Handler() http.Handler {
	return g.engine
}

func (g *GinRouter) ServeStatic(url, absPath string) {
	g.engine.Static(url, absPath)
}
