// Package replication provides cross-region document synchronization.
//
// This file implements the NATS JetStream-backed Replicator.
package replication

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	// defaultStreamRetention is the max duration to keep messages in the stream.
	defaultStreamRetention = 24 * time.Hour

	// defaultMaxBytes is the maximum size of the stream (1GB).
	defaultMaxBytes = 1024 * 1024 * 1024

	// defaultFetchBatch is the number of messages to fetch in one batch.
	defaultFetchBatch = 10

	// defaultFetchTimeout is the max wait time for a fetch request.
	defaultFetchTimeout = 5 * time.Second
)

// NATSReplicator implements Replicator using NATS JetStream.
type NATSReplicator struct {
	nc       *nats.Conn
	js       jetstream.JetStream
	stream   jetstream.Stream
	consumer jetstream.Consumer
	config   Config
	logger   *slog.Logger

	mu      sync.RWMutex
	peers   map[string]Replica
	healthy bool

	// cancelFunc is used to stop the subscription goroutine
	cancelFunc context.CancelFunc
	wg         sync.WaitGroup
}

// NewNATSReplicator creates a new NATS JetStream-backed replicator.
func NewNATSReplicator(cfg Config) (*NATSReplicator, error) {
	if cfg.ServerID == "" {
		return nil, fmt.Errorf("ServerID is required")
	}
	if cfg.Region == "" {
		return nil, fmt.Errorf("Region is required")
	}
	if len(cfg.NATSURLs) == 0 {
		return nil, fmt.Errorf("at least one NATS URL is required")
	}

	// Apply defaults
	if cfg.StreamName == "" {
		cfg.StreamName = DefaultConfig().StreamName
	}
	if cfg.ConsumerDurableName == "" {
		cfg.ConsumerDurableName = fmt.Sprintf("etherply-%s", cfg.ServerID)
	}
	if cfg.ReconnectWait == 0 {
		cfg.ReconnectWait = DefaultConfig().ReconnectWait
	}
	if cfg.MaxReconnects == 0 {
		cfg.MaxReconnects = DefaultConfig().MaxReconnects
	}

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))

	// Connect to NATS with reconnection settings
	opts := []nats.Option{
		nats.Name(cfg.ServerID),
		nats.ReconnectWait(cfg.ReconnectWait),
		nats.MaxReconnects(cfg.MaxReconnects),
		nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
			logger.Warn("nats_disconnected", slog.Any("error", err))
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			logger.Info("nats_reconnected", slog.String("url", nc.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			logger.Info("nats_connection_closed")
		}),
	}

	nc, err := nats.Connect(cfg.NATSURLs[0], opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	// Create JetStream context
	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	ctx := context.Background()

	// Create or get the replication stream
	stream, err := js.CreateOrUpdateStream(ctx, jetstream.StreamConfig{
		Name:        cfg.StreamName,
		Description: "EtherPly cross-region replication events",
		Subjects:    []string{fmt.Sprintf("%s.>", cfg.StreamName)},
		Retention:   jetstream.LimitsPolicy,
		MaxAge:      defaultStreamRetention,
		MaxBytes:    defaultMaxBytes,
		Replicas:    1, // Single replica for now; increase for HA
		Storage:     jetstream.FileStorage,
	})
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("failed to create/update stream: %w", err)
	}

	r := &NATSReplicator{
		nc:      nc,
		js:      js,
		stream:  stream,
		config:  cfg,
		logger:  logger,
		peers:   make(map[string]Replica),
		healthy: true,
	}

	// Register self as a peer
	r.peers[cfg.ServerID] = Replica{
		ID:       cfg.ServerID,
		Region:   cfg.Region,
		Endpoint: cfg.NATSURLs[0],
		LastSeen: time.Now(),
	}

	logger.Info("nats_replicator_initialized",
		slog.String("server_id", cfg.ServerID),
		slog.String("region", cfg.Region),
		slog.String("stream", cfg.StreamName),
	)

	return r, nil
}

// Broadcast sends document changes to all peer replicas.
func (r *NATSReplicator) Broadcast(ctx context.Context, event ChangeEvent) error {
	if !r.Healthy() {
		return fmt.Errorf("replicator is not healthy")
	}

	// Construct subject: STREAM_NAME.changes.<workspace_id>
	subject := fmt.Sprintf("%s.changes.%s", r.config.StreamName, event.WorkspaceID)

	// Serialize the event
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal change event: %w", err)
	}

	// Publish with acknowledgment
	_, err = r.js.Publish(ctx, subject, data)
	if err != nil {
		r.logger.Error("broadcast_failed",
			slog.String("workspace_id", event.WorkspaceID),
			slog.Any("error", err),
		)
		return fmt.Errorf("failed to publish change event: %w", err)
	}

	r.logger.Debug("broadcast_success",
		slog.String("workspace_id", event.WorkspaceID),
		slog.String("subject", subject),
		slog.Int("size_bytes", len(data)),
	)

	return nil
}

// Subscribe registers a handler for incoming changes from peer replicas.
func (r *NATSReplicator) Subscribe(handler ChangeHandler) error {
	ctx, cancel := context.WithCancel(context.Background())
	r.cancelFunc = cancel

	// Create a durable consumer for this server
	consumer, err := r.stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       r.config.ConsumerDurableName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		DeliverPolicy: jetstream.DeliverNewPolicy,
		FilterSubject: fmt.Sprintf("%s.changes.>", r.config.StreamName),
	})
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}
	r.consumer = consumer

	// Start consuming messages in a goroutine
	r.wg.Add(1)
	go func() {
		defer r.wg.Done()
		r.consumeLoop(ctx, handler)
	}()

	r.logger.Info("subscription_started",
		slog.String("consumer", r.config.ConsumerDurableName),
	)

	return nil
}

func (r *NATSReplicator) consumeLoop(ctx context.Context, handler ChangeHandler) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Fetch messages with a timeout
		msgs, err := r.consumer.Fetch(defaultFetchBatch, jetstream.FetchMaxWait(defaultFetchTimeout))
		if err != nil {
			if ctx.Err() != nil {
				return // Context cancelled, exit gracefully
			}
			r.logger.Warn("fetch_error", slog.Any("error", err))
			continue
		}

		for msg := range msgs.Messages() {
			var event ChangeEvent
			if err := json.Unmarshal(msg.Data(), &event); err != nil {
				r.logger.Error("failed_to_unmarshal_event",
					slog.Any("error", err),
					slog.String("subject", msg.Subject()),
				)
				// Ack anyway to skip bad messages
				msg.Ack()
				continue
			}

			// Skip our own messages to prevent loops
			if event.OriginServerID == r.config.ServerID {
				msg.Ack()
				continue
			}

			// Call the handler
			if err := handler(event); err != nil {
				r.logger.Error("handler_error",
					slog.String("workspace_id", event.WorkspaceID),
					slog.Any("error", err),
				)
				// NAK to retry later
				msg.Nak()
				continue
			}

			// Acknowledge successful processing
			msg.Ack()

			r.logger.Debug("event_processed",
				slog.String("workspace_id", event.WorkspaceID),
				slog.String("origin_region", event.OriginRegion),
			)
		}

		if err := msgs.Error(); err != nil {
			r.logger.Warn("message_iteration_error", slog.Any("error", err))
		}
	}
}

// Peers returns the list of known replicas in the cluster.
func (r *NATSReplicator) Peers() []Replica {
	r.mu.RLock()
	defer r.mu.RUnlock()

	peers := make([]Replica, 0, len(r.peers))
	for _, p := range r.peers {
		peers = append(peers, p)
	}
	return peers
}

// Healthy returns true if the replicator is connected and operational.
func (r *NATSReplicator) Healthy() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.healthy && r.nc.IsConnected()
}

// Close shuts down replication connections gracefully.
func (r *NATSReplicator) Close() error {
	r.mu.Lock()
	r.healthy = false
	r.mu.Unlock()

	// Cancel the subscription goroutine
	if r.cancelFunc != nil {
		r.cancelFunc()
	}

	// Wait for goroutines to finish
	r.wg.Wait()

	// Close NATS connection
	r.nc.Close()

	r.logger.Info("nats_replicator_closed")
	return nil
}

// Ensure interface compliance
var _ Replicator = (*NATSReplicator)(nil)
