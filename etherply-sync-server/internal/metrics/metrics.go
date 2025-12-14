package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ConnectedClients tracks the number of currently connected WebSocket clients.
	ConnectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "etherply_connections_active",
		Help: "The total number of currently active WebSocket connections",
	})

	// MessagesReceived tracks the total number of messages received from clients.
	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "etherply_messages_received_total",
		Help: "The total number of messages received from clients",
	})

	// MessagesBroadcast tracks the total number of messages broadcast to clients.
	MessagesBroadcast = promauto.NewCounter(prometheus.CounterOpts{
		Name: "etherply_messages_broadcast_total",
		Help: "The total number of messages broadcast to clients",
	})

	// OperationsTotal tracks the total number of CRDT operations processed.
	OperationsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "etherply_operations_total",
		Help: "The total number of processed CRDT operations",
	})

	// OperationDuration tracks the time taken to process CRDT operations.
	OperationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "etherply_operation_duration_seconds",
		Help:    "Time taken to process CRDT operations in seconds",
		Buckets: prometheus.DefBuckets,
	})

	// SyncLatency tracks the time taken to process sync messages (generic).
	SyncLatency = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "etherply_sync_latency_seconds",
		Help:    "Time taken to process sync messages in seconds",
		Buckets: prometheus.DefBuckets,
	})
)
