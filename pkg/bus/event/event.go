package bus

import "context"

// Event represents a domain event.
type Event interface {
	EventName() string
}

// EventHandler defines the contract for handling events.
type EventHandler interface {
	Handle(ctx context.Context, event Event) error
}

// EventBus is the interface for subscribing to and publishing events.
type EventBus interface {
	Publish(ctx context.Context, event Event) error
}
