package command

import (
	"context"
	"github.com/mik3lon/starter-template/pkg/bus"
)

type CommandHandler interface {
	Handle(ctx context.Context, command bus.Dto) error
}
