package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// ConnectedClients tracks the number of currently connected WebSocket clients
	ConnectedClients = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "etherply_connected_clients",
		Help: "The total number of connected clients",
	})

	// MessagesReceived tracks total messages received from clients
	MessagesReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "etherply_messages_received_total",
		Help: "The total number of messages received from clients",
	})

	// MessagesBroadcast tracks total messages sent to clients
	MessagesBroadcast = promauto.NewCounter(prometheus.CounterOpts{
		Name: "etherply_messages_broadcast_total",
		Help: "The total number of messages broadcast to clients",
	})

	// OperationDuration tracks the time taken to process operations
	OperationDuration = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "etherply_operation_duration_seconds",
		Help:    "Time taken to process sync operations",
		Buckets: prometheus.DefBuckets,
	})
)
