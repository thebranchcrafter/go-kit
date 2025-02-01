package application_event

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
)

// EventBus is the interface for subscribing to and publishing events.
type EventBus interface {
	Publish(ctx context.Context, event domain.Event) error
}
