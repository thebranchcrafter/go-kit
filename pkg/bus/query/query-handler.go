package query

import (
	"context"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type QueryHandler interface {
	Handle(ctx context.Context, query bus.Dto) (interface{}, error)
}
