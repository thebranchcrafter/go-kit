package domain

// AggregateRoot defines the base structure for an aggregate.
type AggregateRoot interface {
	Id() string                // Returns the unique identifier
	Record(event Event)        // Stores a domain event
	PullDomainEvents() []Event // Retrieves and clears stored domain events
}
