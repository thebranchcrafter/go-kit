package application_command

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/application"
)

type CommandHandler interface {
	Handle(ctx context.Context, c application.Command) error
}
