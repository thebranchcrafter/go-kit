package infrastructure

import (
	"context"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

// NatsBroker implements the domain.Broker interface
type NatsBroker struct {
	conn   *nats.Conn
	sub    *nats.Subscription
	msgCh  chan *nats.Msg
	mu     sync.Mutex
	closed bool
}

// NewNatsBroker creates a new NATS broker connection
func NewNatsBroker(url, subject string) (*NatsBroker, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	msgCh := make(chan *nats.Msg, 64) // Buffered channel for message processing

	sub, err := nc.ChanSubscribe(subject, msgCh)
	if err != nil {
		nc.Close()
		return nil, err
	}

	return &NatsBroker{
		conn:  nc,
		sub:   sub,
		msgCh: msgCh,
	}, nil
}

// FetchMessage fetches a message from the NATS subject
func (n *NatsBroker) FetchMessage(_ context.Context) ([]byte, error) {
	msg := <-n.msgCh // Blocks until a message is received
	return msg.Data, nil
}

// Close gracefully shuts down the NATS connection
func (n *NatsBroker) Close() {
	n.mu.Lock()
	defer n.mu.Unlock()
	if !n.closed {
		_ = n.sub.Unsubscribe()
		n.conn.Close()
		n.closed = true
		log.Println("NATS connection closed")
	}
}
