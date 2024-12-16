package mux_router

import (
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/router"
	"net/http"

	"github.com/gorilla/mux"
)

type MuxRouter struct {
	router *mux.Router
}

func NewMuxRouter() *MuxRouter {
	return &MuxRouter{
		router: mux.NewRouter(),
	}
}

func (m *MuxRouter) Handle(method, path string, handler http.HandlerFunc, middleware ...router.Middleware) {
	// Apply middleware
	finalHandler := handler
	for _, mware := range middleware {
		finalHandler = mware(finalHandler)
	}

	// Register the route
	m.router.HandleFunc(path, finalHandler).Methods(method)
}

func (m *MuxRouter) Serve(addr string) error {
	return http.ListenAndServe(addr, m.router)
}

func (m *MuxRouter) Handler() http.Handler {
	return m.router
}

func (m *MuxRouter) ServeStatic(url, absPath string) http.Handler {
	fileServer := http.FileServer(http.Dir(absPath))
	m.router.PathPrefix(url).Handler(http.StripPrefix(url, fileServer))

	return fileServer
}
