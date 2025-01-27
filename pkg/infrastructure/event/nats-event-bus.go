package infrastructure_event

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
	"time"

	"github.com/nats-io/nats.go"
)

// NATSEventBus is an implementation of EventBus using NATS.
type NATSEventBus struct {
	conn *nats.Conn
}

// NewNATSEventBus creates a new instance of NATSEventBus with reconnection options.
func NewNATSEventBus(url string) (*NATSEventBus, error) {
	conn, err := nats.Connect(
		url,
		nats.MaxReconnects(-1),            // Unlimited reconnection attempts
		nats.ReconnectWait(2*time.Second), // Wait 2 seconds between reconnection attempts
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			fmt.Printf("NATS disconnected: %v\n", err)
		}),
		nats.ReconnectHandler(func(conn *nats.Conn) {
			fmt.Printf("NATS reconnected to %s\n", conn.ConnectedUrl())
		}),
		nats.ClosedHandler(func(conn *nats.Conn) {
			fmt.Println("NATS connection closed")
		}),
		nats.PingInterval(10*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}
	return &NATSEventBus{conn: conn}, nil
}

// Publish publishes an event to a NATS subject.
func (b *NATSEventBus) Publish(_ context.Context, event domain.Event) error {
	// Serialize the event to JSON
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Publish the event to a NATS subject based on the event name
	subject := event.EventName()
	if err := b.conn.Publish(subject, data); err != nil {
		return fmt.Errorf("failed to publish event to NATS: %w", err)
	}

	// Ensure the message is flushed to the server
	if err := b.conn.Flush(); err != nil {
		return fmt.Errorf("failed to flush NATS connection: %w", err)
	}

	// Check for errors during the flush
	if err := b.conn.LastError(); err != nil {
		return fmt.Errorf("NATS connection error: %w", err)
	}

	return nil
}

// Close closes the NATS connection.
func (b *NATSEventBus) Close() {
	b.conn.Close()
}
