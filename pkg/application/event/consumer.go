package application_event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
	"log"
)

// EventConsumer represents a generic consumer.
type EventConsumer struct {
	broker       domain.Broker
	event        domain.Event
	handler      domain.EventHandler
	messageName  string
	errorChannel chan ErrorMessage
}

// ErrorMessage represents an error and its associated message.
type ErrorMessage struct {
	Error error
	Msg   []byte
}

// NewEventConsumer creates a new EventConsumer with an error channel.
func NewEventConsumer(
	broker domain.Broker,
	event domain.Event,
	handler domain.EventHandler,
	messageName string,
	errorChannel chan ErrorMessage,
) *EventConsumer {
	return &EventConsumer{
		broker:       broker,
		event:        event,
		handler:      handler,
		messageName:  messageName,
		errorChannel: errorChannel,
	}
}

// Start starts the consumer and processes messages until a stop signal is received.
func (c *EventConsumer) Start(ctx context.Context, stopChan chan struct{}) {
	log.Printf("Starting consumer for: %s\n", c.messageName)

	for {
		select {
		case <-stopChan:
			log.Printf("Stopping consumer for: %s\n", c.messageName)
			return
		default:
			func() {
				// Recover from panic for this message processing
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Recovered from panic: %v", r)
						c.sendError(fmt.Errorf("panic: %v", r), nil)
					}
				}()

				// Fetch message
				msg, err := c.broker.FetchMessage(ctx)
				if err != nil {
					log.Printf("Error fetching message: %s\n", err)
					c.sendError(err, nil)
					return
				}

				// Deserialize message
				var payload map[string]interface{}
				if err := json.Unmarshal(msg, &payload); err != nil {
					log.Printf("Error unmarshalling message: %s\n", err)
					c.sendError(err, msg)
					return
				}

				// Map to domain event
				if err := c.event.FromMap(payload); err != nil {
					log.Printf("Error building domain event: %s\n", err)
					c.sendError(err, msg)
					return
				}

				// Handle the event
				if err := c.handler.Handle(ctx, c.event); err != nil {
					log.Printf("Error processing message: %s\n", err)
					c.sendError(err, msg)
					return
				}
			}()
		}
	}
}

// sendError sends the error and message to the error channel.
func (c *EventConsumer) sendError(err error, msg []byte) {
	if c.errorChannel != nil {
		c.errorChannel <- ErrorMessage{Error: err, Msg: msg}
	}
}
