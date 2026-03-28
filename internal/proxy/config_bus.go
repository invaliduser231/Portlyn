package proxy

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type InMemoryConfigBus struct {
	mu          sync.RWMutex
	subscribers map[chan RouteChangedEvent]struct{}
}

func NewInMemoryConfigBus() *InMemoryConfigBus {
	return &InMemoryConfigBus{subscribers: make(map[chan RouteChangedEvent]struct{})}
}

func (b *InMemoryConfigBus) PublishRouteChanged(_ context.Context, host string) error {
	event := RouteChangedEvent{
		Host:      normalizeHost(host),
		ChangedAt: time.Now().UTC(),
	}

	b.mu.RLock()
	defer b.mu.RUnlock()
	for ch := range b.subscribers {
		select {
		case ch <- event:
		default:
		}
	}
	return nil
}

func (b *InMemoryConfigBus) SubscribeRouteChanged(ctx context.Context) <-chan RouteChangedEvent {
	ch := make(chan RouteChangedEvent, 32)

	b.mu.Lock()
	b.subscribers[ch] = struct{}{}
	b.mu.Unlock()

	go func() {
		<-ctx.Done()
		b.mu.Lock()
		delete(b.subscribers, ch)
		close(ch)
		b.mu.Unlock()
	}()

	return ch
}

type RedisConfigBus struct {
	client  *redis.Client
	channel string
}

func NewRedisConfigBus(client *redis.Client, channel string) *RedisConfigBus {
	if channel == "" {
		channel = "portlyn:proxy:route-changed"
	}
	return &RedisConfigBus{
		client:  client,
		channel: channel,
	}
}

func (b *RedisConfigBus) PublishRouteChanged(ctx context.Context, host string) error {
	payload, err := json.Marshal(RouteChangedEvent{
		Host:      normalizeHost(host),
		ChangedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	return b.client.Publish(ctx, b.channel, payload).Err()
}

func (b *RedisConfigBus) SubscribeRouteChanged(ctx context.Context) <-chan RouteChangedEvent {
	ch := make(chan RouteChangedEvent, 32)
	pubsub := b.client.Subscribe(ctx, b.channel)

	go func() {
		defer close(ch)
		defer pubsub.Close()

		msgCh := pubsub.Channel()
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgCh:
				if !ok {
					return
				}
				var event RouteChangedEvent
				if err := json.Unmarshal([]byte(msg.Payload), &event); err != nil {
					continue
				}
				event.Host = normalizeHost(event.Host)
				select {
				case ch <- event:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return ch
}
