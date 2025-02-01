package domain

type EventRecorder struct {
	events []Event
}

func (d *EventRecorder) AddDomainEvent(event Event) {
	d.events = append(d.events, event)
}

func (d *EventRecorder) PullDomainEvents() []Event {
	return d.events
}
