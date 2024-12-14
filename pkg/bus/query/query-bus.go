package query

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/bus"
	"github.com/thebranchcrafter/go-kit/pkg/infrastructure/logger"
	"sync"
)

type Bus interface {
	RegisterQuery(query bus.Query, handler QueryHandler) error
	Ask(ctx context.Context, q bus.Query) (interface{}, error)
}

type QueryBus struct {
	handlers map[string]QueryHandler
	lock     sync.Mutex
	logger   logger.Logger
}

func InitQueryBus(l logger.Logger) *QueryBus {
	return &QueryBus{
		handlers: make(map[string]QueryHandler, 0),
		lock:     sync.Mutex{},
		logger:   l,
	}
}

type QueryAlreadyRegistered struct {
	message   string
	queryName string
}

func (i QueryAlreadyRegistered) Error() string {
	return i.message
}

func NewQueryAlreadyRegistered(message string, queryName string) QueryAlreadyRegistered {
	return QueryAlreadyRegistered{message: message, queryName: queryName}
}

type QueryNotRegistered struct {
	message   string
	queryName string
}

func (i QueryNotRegistered) Error() string {
	return i.message
}

func NewQueryNotRegistered(message string, queryName string) QueryAlreadyRegistered {
	return QueryAlreadyRegistered{message: message, queryName: queryName}
}

func (bus *QueryBus) RegisterQuery(query bus.Query, handler QueryHandler) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	queryName := query.Id()

	if _, ok := bus.handlers[queryName]; ok {
		return NewQueryAlreadyRegistered("Query already registered", queryName)
	}

	bus.handlers[queryName] = handler

	return nil
}

func (bus *QueryBus) Ask(ctx context.Context, query bus.Query) (interface{}, error) {
	queryName := query.Id()

	if handler, ok := bus.handlers[queryName]; ok {
		response, err := bus.doAsk(ctx, handler, query)
		if err != nil {
			return nil, err
		}

		return response, nil
	}

	return nil, NewQueryNotRegistered("Query not registered", queryName)
}

func (bus *QueryBus) doAsk(ctx context.Context, handler QueryHandler, query bus.Query) (interface{}, error) {
	return handler.Handle(ctx, query)
}

type QueryNotValid struct {
	message string
}

func (i QueryNotValid) Error() string {
	return i.message
}
