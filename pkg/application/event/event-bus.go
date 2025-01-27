package application_event

import (
	"context"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
)

// Bus is the interface for subscribing to and publishing events.
type Bus interface {
	Publish(ctx context.Context, event domain.Event) error
}
