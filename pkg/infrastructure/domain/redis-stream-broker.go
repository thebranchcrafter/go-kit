package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/thebranchcrafter/go-kit/pkg/domain"
)

// RedisStreamBroker implements domain.Broker and domain.EventBus using Redis Streams.
type RedisStreamBroker struct {
	client     *redis.Client
	streamName string
	groupName  string
	consumerID string
	closed     bool
}

// NewRedisStreamBroker initializes a Redis Stream broker.
func NewRedisStreamBroker(redisAddr, streamName, groupName, consumerID string) (*RedisStreamBroker, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Ensure the stream and consumer group exist
	err := rdb.XGroupCreateMkStream(context.Background(), streamName, groupName, "$").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return nil, fmt.Errorf("failed to create Redis stream group: %w", err)
	}

	return &RedisStreamBroker{
		client:     rdb,
		streamName: streamName,
		groupName:  groupName,
		consumerID: consumerID,
		closed:     false,
	}, nil
}

// Publish sends a domain event to the Redis stream.
func (r *RedisStreamBroker) Publish(ctx context.Context, event domain.Event) error {
	if r.closed {
		return fmt.Errorf("broker is closed")
	}

	// Serialize event to JSON
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return fmt.Errorf("failed to serialize event: %w", err)
	}

	// Publish event to Redis Stream
	_, err = r.client.XAdd(ctx, &redis.XAddArgs{
		Stream: r.streamName,
		Values: map[string]interface{}{
			"aggregate_id":   event.AggregateID(),
			"event_name":     event.EventName(),
			"occurred_at":    event.OccurredOn().Format(time.RFC3339),
			"correlation_id": event.CorrelationID(),
			"payload":        string(payload), // Store payload as JSON string
		},
	}).Result()

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event: %s to stream: %s", event.EventName(), r.streamName)
	return nil
}

// FetchMessage retrieves a message from the Redis Stream.
func (r *RedisStreamBroker) FetchMessage(ctx context.Context) ([]byte, error) {
	if r.closed {
		return nil, fmt.Errorf("broker is closed")
	}

	streams, err := r.client.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    r.groupName,
		Consumer: r.consumerID,
		Streams:  []string{r.streamName, ">"},
		Count:    1,
		Block:    5 * time.Second, // Block waiting for new messages
	}).Result()

	if err == redis.Nil {
		return nil, nil // No new messages
	} else if err != nil {
		return nil, fmt.Errorf("failed to read from Redis stream: %w", err)
	}

	if len(streams) > 0 && len(streams[0].Messages) > 0 {
		msg := streams[0].Messages[0]

		// Acknowledge message
		_ = r.client.XAck(context.Background(), r.streamName, r.groupName, msg.ID).Err()

		// Convert message to JSON
		data, err := json.Marshal(msg.Values)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize message: %w", err)
		}

		return data, nil
	}

	return nil, nil
}

// Close shuts down the Redis broker.
func (r *RedisStreamBroker) Close() {
	if r.closed {
		return
	}
	r.closed = true
	_ = r.client.Close()
	log.Println("Redis stream broker closed.")
}
