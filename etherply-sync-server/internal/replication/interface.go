// Package replication provides cross-region document synchronization
// for Multi-Region Replication support in EtherPly.
//
// This package enables Active-Active replication where clients can write
// to any region and documents converge automatically via Automerge CRDTs.
package replication

import (
	"context"
	"time"
)

// Replica represents a peer node in the replication cluster.
type Replica struct {
	// ID is a unique identifier for this replica (e.g., "sync-us-east-1").
	ID string `json:"id"`

	// Endpoint is the connection string (e.g., "nats://us-east-1.etherply.io:4222").
	Endpoint string `json:"endpoint"`

	// Region is the geographic region identifier (e.g., "us-east-1", "eu-west-1").
	Region string `json:"region"`

	// LastSeen is the timestamp of the last heartbeat from this replica.
	LastSeen time.Time `json:"last_seen"`
}

// ChangeEvent represents a document change to be replicated.
type ChangeEvent struct {
	// WorkspaceID identifies the workspace/document being changed.
	WorkspaceID string `json:"workspace_id"`

	// Changes is the serialized Automerge change data.
	Changes []byte `json:"changes"`

	// OriginRegion identifies where this change originated (to prevent loops).
	OriginRegion string `json:"origin_region"`

	// OriginServerID identifies the specific server that produced this change.
	OriginServerID string `json:"origin_server_id"`

	// Timestamp is when the change was created.
	Timestamp time.Time `json:"timestamp"`
}

// ChangeHandler is a callback function invoked when changes are received from peers.
type ChangeHandler func(event ChangeEvent) error

// Replicator manages cross-region document synchronization.
// Implementations must be thread-safe.
type Replicator interface {
	// Broadcast sends document changes to all peer replicas.
	// This is a fire-and-forget operation; delivery is best-effort but durable
	// (changes are persisted in the message queue for later delivery).
	//
	// Returns an error if the broadcast could not be queued (e.g., connection failure).
	Broadcast(ctx context.Context, event ChangeEvent) error

	// Subscribe registers a handler for incoming changes from peer replicas.
	// The handler is called for each change event received.
	//
	// Returns an error if subscription could not be established.
	Subscribe(handler ChangeHandler) error

	// Peers returns the list of known replicas in the cluster.
	// This includes healthy and potentially unhealthy peers.
	Peers() []Replica

	// Healthy returns true if the replicator is connected and operational.
	Healthy() bool

	// Close shuts down replication connections gracefully.
	// It waits for in-flight messages to be acknowledged before returning.
	Close() error
}

// Config holds configuration for a Replicator.
type Config struct {
	// ServerID is a unique identifier for this server instance.
	ServerID string

	// Region is the geographic region this server is deployed in.
	Region string

	// NATSURLs is a list of NATS server URLs to connect to.
	NATSURLs []string

	// StreamName is the NATS JetStream stream name for replication events.
	// Default: "ETHERPLY_REPLICATION"
	StreamName string

	// ConsumerDurableName is used for durable NATS consumers.
	// Default: derived from ServerID
	ConsumerDurableName string

	// ReconnectWait is the duration to wait between reconnection attempts.
	// Default: 2 seconds
	ReconnectWait time.Duration

	// MaxReconnects is the maximum number of reconnection attempts.
	// Set to -1 for unlimited. Default: -1
	MaxReconnects int
}

// DefaultConfig returns a Config with sensible defaults filled in.
func DefaultConfig() Config {
	return Config{
		StreamName:    "ETHERPLY_REPLICATION",
		ReconnectWait: 2 * time.Second,
		MaxReconnects: -1,
	}
}
