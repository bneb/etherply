package pubsub

import (
	"sync"

	"github.com/google/uuid"
)

// MemoryPubSub is a thread-safe in-memory implementation of PubSub.
type MemoryPubSub struct {
	mu     sync.RWMutex
	topics map[string]map[string]chan Message // topic -> subID -> chan
}

func NewMemoryPubSub() *MemoryPubSub {
	return &MemoryPubSub{
		topics: make(map[string]map[string]chan Message),
	}
}

func (ps *MemoryPubSub) Publish(topic string, msg Message) error {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	subs, ok := ps.topics[topic]
	if !ok {
		return nil // No subscribers
	}

	for _, ch := range subs {
		// Non-blocking send to prevent one slow subscriber from blocking publisher?
		// Or blocking? For in-memory Sync Server, likely blocking with buffer is safer to assume delivery,
		// but if buffer full, we drop or block?
		// Let's use non-blocking select to avoid deadlock if subscriber is dead.
		select {
		case ch <- msg:
		default:
			// Buffer full, drop message or log?
			// For reliable sync, this is bad. But if client is too slow, maybe we disconnect them.
			// Ideally we use a larger buffer.
		}
	}
	return nil
}

func (ps *MemoryPubSub) Subscribe(topic string) (<-chan Message, func()) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if _, ok := ps.topics[topic]; !ok {
		ps.topics[topic] = make(map[string]chan Message)
	}

	// Buffer size 100 to tolerate bursts
	ch := make(chan Message, 100)
	subID := uuid.New().String()

	ps.topics[topic][subID] = ch

	unsub := func() {
		ps.mu.Lock()
		defer ps.mu.Unlock()

		if subs, ok := ps.topics[topic]; ok {
			if _, exists := subs[subID]; exists {
				delete(subs, subID)
				close(ch)
			}
			if len(subs) == 0 {
				delete(ps.topics, topic)
			}
		}
	}

	return ch, unsub
}

func (ps *MemoryPubSub) Stats() PubSubStats {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	activeTopics := len(ps.topics)
	totalSubs := 0
	for _, subs := range ps.topics {
		totalSubs += len(subs)
	}

	return PubSubStats{
		ActiveTopics:     activeTopics,
		TotalSubscribers: totalSubs,
	}
}
