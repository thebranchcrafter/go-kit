package application_event

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
)

// EventHandler defines the contract for handling events.
type EventHandler interface {
	Handle(ctx context.Context, event domain.Event) error
}
