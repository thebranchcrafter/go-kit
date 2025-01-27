package application_query

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/application"
)

type QueryHandler interface {
	Handle(ctx context.Context, query application.Query) (interface{}, error)
}
