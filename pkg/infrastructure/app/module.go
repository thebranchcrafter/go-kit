package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/thebranchcrafter/go-kit/pkg/application"
	application_command "github.com/thebranchcrafter/go-kit/pkg/application/command"
	application_query "github.com/thebranchcrafter/go-kit/pkg/application/query"
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/router"
)

type Modules []Module

// Route represents an HTTP route with its method, path, handler, and middleware.
type Route struct {
	Method      string
	Path        string
	Handler     gin.HandlerFunc
	Middlewares []router.Middleware
}

// Module represents a module that can register routes.
type Module interface {
	Routes() []Route
	Name() string
	Commands() map[application.Command]application_command.CommandHandler
	Queries() map[application.Query]application_query.QueryHandler
}

type BaseModule struct {
	commands map[application.Command]application_command.CommandHandler
	queries  map[application.Query]application_query.QueryHandler
	CommonDependencies
}

// AddCommand adds a command to the module
func (bm *BaseModule) AddCommand(c application.Command, commandHandler application_command.CommandHandler) {
	if bm.commands == nil {
		bm.commands = make(map[application.Command]application_command.CommandHandler)
	}
	bm.commands[c] = commandHandler
}

// AddQuery adds a query to the module
func (bm *BaseModule) AddQuery(c application.Query, queryHandler application_query.QueryHandler) {
	if bm.queries == nil {
		bm.queries = make(map[application.Query]application_query.QueryHandler)
	}
	bm.queries[c] = queryHandler
}

// Commands returns all commands registered in the module
func (bm *BaseModule) Commands() map[application.Command]application_command.CommandHandler {
	return bm.commands
}

// Queries returns all commands registered in the module
func (bm *BaseModule) Queries() map[application.Query]application_query.QueryHandler {
	return bm.queries
}

type AlreadyExistsError struct {
	m Module
}

func NewModuleAlreadyExistsError(m Module) *AlreadyExistsError {
	return &AlreadyExistsError{m: m}
}

func (m AlreadyExistsError) Error() string {
	return fmt.Sprintf("module %s already exists", m.m.Name())
}
