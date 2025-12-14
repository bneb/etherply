package pubsub

// Message represents a broadcast message in a topic.
type Message struct {
	Topic    string
	Payload  []byte
	SenderID string // Optional: exclude this sender from receiving (logic depends on implementation)
}

// PubSubStats contains metrics about the PubSub system.
type PubSubStats struct {
	TotalSubscribers int `json:"total_subscribers"`
	ActiveTopics     int `json:"active_topics"`
}

// PubSub defines the behavior for a publish-subscribe system.
type PubSub interface {
	// Publish broadcasts a message to a topic.
	Publish(topic string, msg Message) error

	// Subscribe joins a topic and returns a read-only channel for messages and an unsubscribe function.
	Subscribe(topic string) (<-chan Message, func())

	// Stats returns current metrics.
	Stats() PubSubStats
}
