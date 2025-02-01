package domain

import "context"

type Broker interface {
	FetchMessage(ctx context.Context) ([]byte, error)
	Close()
}
