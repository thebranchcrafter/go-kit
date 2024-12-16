package query

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/bus"
)

type QueryHandler interface {
	Handle(ctx context.Context, query bus.Query) (interface{}, error)
}
