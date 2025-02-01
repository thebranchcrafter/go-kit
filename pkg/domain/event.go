package domain

import "time"

// Event represents a domain event.
type Event interface {
	AggregateID() string
	OccurredOn() time.Time
	EventName() string
	Payload() map[string]interface{}
	Version() int
	CorrelationID() string
	FromMap(data map[string]interface{}) error
}
