package domain

import (
	"context"
)

// EventHandler defines the contract for handling events.
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}
