package event

import "context"

// EventBus is the interface for subscribing to and publishing events.
type EventBus interface {
	Publish(ctx context.Context, event Event) error
}
