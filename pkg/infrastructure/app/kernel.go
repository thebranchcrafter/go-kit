package app

import (
	"context"
	"github.com/gin-gonic/gin"
	application_command "github.com/thebranchcrafter/go-kit/pkg/application/command"
	application_event "github.com/thebranchcrafter/go-kit/pkg/application/event"
	application_query "github.com/thebranchcrafter/go-kit/pkg/application/query"
	"net/http"

	http_response "github.com/thebranchcrafter/go-kit/pkg/infrastructure/http/response"
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/logger"
)

type CommonDependencies struct {
	Router         *gin.Engine
	CommandBus     application_command.Bus
	QueryBus       application_query.Bus
	EventBus       application_event.Bus
	ResponseWriter http_response.ResponseWriter
	Logger         logger.Logger
}

// Kernel holds the core infrastructure and components.
type Kernel struct {
	server  *http.Server
	Modules map[string]Module
	CommonDependencies
}

// NewKernel creates a new Kernel instance with functional options.
func NewKernel(options ...func(*Kernel)) *Kernel {
	k := &Kernel{
		Modules: make(map[string]Module),
	}
	for _, opt := range options {
		opt(k)
	}
	return k
}

// WithRouter sets a custom router implementation.
func WithRouter(r *gin.Engine) func(*Kernel) {
	return func(k *Kernel) {
		k.Router = r
		if k.server != nil {
			k.server.Handler = r.Handler()
		}
	}
}

// WithCommandBus sets a custom CommandBus.
func WithCommandBus(cb application_command.Bus) func(*Kernel) {
	return func(k *Kernel) {
		k.CommandBus = cb
	}
}

// WithEventBus sets a custom WithEventBus.
func WithEventBus(eb application_event.Bus) func(*Kernel) {
	return func(k *Kernel) {
		k.EventBus = eb
	}
}

// WithQueryBus sets a custom QueryBus.
func WithQueryBus(qb application_query.Bus) func(*Kernel) {
	return func(k *Kernel) {
		k.QueryBus = qb
	}
}

// WithLogger sets a custom logger implementation.
func WithLogger(l logger.Logger) func(*Kernel) {
	return func(k *Kernel) {
		k.Logger = l
	}
}

// WithJsonResponseWriter sets a custom JSON response writer.
func WithJsonResponseWriter(w http_response.ResponseWriter) func(*Kernel) {
	return func(k *Kernel) {
		k.ResponseWriter = w
	}
}

// AddModule adds a module to the kernel.
func (k *Kernel) AddModule(m Module) error {
	if k.Modules == nil {
		k.Modules = make(map[string]Module)
	}

	if k.Modules[m.Name()] != nil {
		return NewModuleAlreadyExistsError(m)
	}
	k.Modules[m.Name()] = m

	for c, ch := range m.Commands() {
		if err := k.CommandBus.RegisterCommand(c, ch); err != nil {
			return err
		}
	}

	for q, ch := range m.Queries() {
		if err := k.QueryBus.RegisterQuery(q, ch); err != nil {
			return err
		}
	}

	return nil
}

// RegisterRoutes allows each module to register its routes.
func (k *Kernel) RegisterRoutes() {
	for _, module := range k.Modules {
		for _, route := range module.Routes() {
			// Apply middleware if any
			handler := route.Handler
			for _, mw := range route.Middlewares {
				handler = mw(handler)
			}
			// Register the route in the router
			k.Router.Handle(route.Method, route.Path, handler)
		}
	}
}

// StartServer starts the HTTP server.
func (k *Kernel) StartServer(port string) error {
	k.server = &http.Server{
		Addr:    port,
		Handler: k.Router.Handler(),
	}

	return k.server.ListenAndServe()
}

// ShutdownServer gracefully shuts down the HTTP server.
func (k *Kernel) ShutdownServer(ctx context.Context) error {
	return k.server.Shutdown(ctx)
}
