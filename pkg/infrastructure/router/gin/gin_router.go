package gin_router

import (
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/router"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	engine *gin.Engine
}

func NewGinRouter() *GinRouter {
	// gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	return &GinRouter{
		engine: engine,
	}
}

func (g *GinRouter) Handle(method, path string, handler http.HandlerFunc, middleware ...router.Middleware) {
	// Wrap the http.HandlerFunc into a gin.HandlerFunc
	ginHandler := gin.WrapH(handler)

	for _, m := range middleware {
		ginHandler = wrapMiddlewareGin(m, ginHandler)
	}

	// Register the route
	switch method {
	case http.MethodGet:
		g.engine.GET(path, ginHandler)
	case http.MethodPost:
		g.engine.POST(path, ginHandler)
	case http.MethodPut:
		g.engine.PUT(path, ginHandler)
	case http.MethodDelete:
		g.engine.DELETE(path, ginHandler)
	}
}

func (g *GinRouter) Serve(addr string) error {
	return g.engine.Run(addr)
}

func (g *GinRouter) Handler() http.Handler {
	return g.engine
}

// Wrap middleware from router.Middleware to gin.HandlerFunc
func wrapMiddlewareGin(m router.Middleware, next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		m(func(w http.ResponseWriter, r *http.Request) {
			next(c)
		})(c.Writer, c.Request)
	}
}
