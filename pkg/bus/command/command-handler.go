package command

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/bus"
)

type CommandHandler interface {
	Handle(ctx context.Context, c bus.Command) error
}
