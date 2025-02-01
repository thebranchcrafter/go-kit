package domain

import "sync"

// AggregateRoot defines the base structure for an aggregate.
type AggregateRoot interface {
	Id() string                 // Returns the unique identifier
	AddDomainEvent(event Event) // Stores a domain event
	PullDomainEvents() []Event  // Retrieves and clears stored domain events
}

// BaseAggregateRoot provides reusable behavior for all aggregate roots.
type BaseAggregateRoot struct {
	id          string     // Aggregate ID
	events      []Event    // Stored domain events
	version     int        // Version for concurrency control
	eventsMutex sync.Mutex // Ensures thread safety when adding/pulling events
}

// NewBaseAggregateRoot initializes a new aggregate root.
func NewBaseAggregateRoot(id string) *BaseAggregateRoot {
	return &BaseAggregateRoot{
		id:      id,
		events:  []Event{},
		version: 1,
	}
}

// Id returns the aggregate's unique identifier.
func (a *BaseAggregateRoot) Id() string {
	return a.id
}

// Version returns the aggregate version (useful for optimistic concurrency).
func (a *BaseAggregateRoot) Version() int {
	return a.version
}

// AddDomainEvent appends a domain event to the aggregate's event list.
func (a *BaseAggregateRoot) AddDomainEvent(event Event) {
	a.eventsMutex.Lock()
	defer a.eventsMutex.Unlock()
	a.events = append(a.events, event)
}

// PullDomainEvents retrieves and clears stored domain events.
func (a *BaseAggregateRoot) PullDomainEvents() []Event {
	a.eventsMutex.Lock()
	defer a.eventsMutex.Unlock()

	// Copy events to return and clear the slice
	pulledEvents := a.events
	a.events = []Event{}

	return pulledEvents
}
