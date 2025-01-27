package domain

// Event represents a domain event.
type Event interface {
	EventName() string
}
