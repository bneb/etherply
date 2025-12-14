package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/bneb/etherply/etherply-sync-server/internal/config"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/bneb/etherply/etherply-sync-server/internal/pubsub"
	"github.com/bneb/etherply/etherply-sync-server/internal/replication"
	"github.com/bneb/etherply/etherply-sync-server/internal/server"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
	"github.com/bneb/etherply/etherply-sync-server/internal/sync"
	"github.com/bneb/etherply/etherply-sync-server/internal/webhook"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// main is the entry point for the EtherPly Sync Server.
//
// Architecture Overview:
//
//  1. Config: Loads configuration from environment variables (12-factor app).
//  2. Persistence (BadgerDB): We use BadgerDB (LSM tree) instead of BoltDB because
//     sync engines are write-heavy. Persistent allows recovery from crash loops.
//  3. Engine: Starts the sync engine with configurable strategy (CRDT, LWW, etc).
//  4. Replication: Optionally enables multi-region replication via NATS JetStream.
//  5. Presence: Starts the ephemeral Presence Manager (Redis-backed in prod, memory in dev).
//  6. HTTP: Sets up routes and health checks (readiness/liveness probes).
//  7. Graceful Shutdown: Handles SIGTERM for K8s rolling updates. A hard kill
//     risks corrupting the LSM tree.
//
// Configuration:
//   - SYNC_STRATEGY: automerge (default), lww, server-auth
//   - ETHERPLY_JWT_SECRET: Required for authentication
//   - BADGER_PATH: Storage path (default: ./badger.db)
//   - NATS_URL: Enable multi-region replication
//   - LOG_FORMAT: json (default), text
//   - LOG_LEVEL: debug, info, warn, error
func main() {
	// Load configuration from environment
	cfg := config.Load()

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		log.Fatalf("[CRITICAL] Configuration error: %v", err)
	}

	// Initialize logger
	logger := cfg.NewLogger()

	// Initialize authentication
	auth.Init(cfg.JWTSecret)

	// Initialize Store (BadgerDB v4)
	// CRITICAL: Ensure this path is mounted on a PVC in Kubernetes.
	// Ephemeral storage will lead to dataloss on pod restart.
	stateStore, err := store.NewBadgerStore(cfg.BadgerPath)
	if err != nil {
		log.Fatalf("Failed to initialize persistence layer at %s: %v", cfg.BadgerPath, err)
	}
	defer stateStore.Close()
	logger.Info("persistence_initialized", "path", cfg.BadgerPath)

	// Create sync strategy based on configuration
	strategy := sync.NewStrategy(cfg.SyncStrategy)
	logger.Info("sync_strategy_selected", "strategy", strategy.Name())

	// Initialize CRDT Engine with configured strategy
	crdtEngine := crdt.NewEngine(stateStore,
		crdt.WithStrategy(strategy),
		crdt.WithLogger(logger),
	)

	// Initialize Multi-Region Replication (if configured)
	var replicator *replication.NATSReplicator
	if len(cfg.NATSURLs) > 0 {
		replicator, err = replication.NewNATSReplicator(replication.Config{
			ServerID: cfg.ServerID,
			Region:   cfg.Region,
			NATSURLs: cfg.NATSURLs,
		})
		if err != nil {
			log.Fatalf("Failed to initialize NATS replicator: %v", err)
		}
		defer replicator.Close()

		// Wire replication to engine
		crdtEngine.SetReplicator(replicator, cfg.Region, cfg.ServerID)

		// Subscribe to incoming changes
		if err := replicator.Subscribe(func(event replication.ChangeEvent) error {
			return crdtEngine.ApplyRemoteChanges(event.WorkspaceID, event.Changes)
		}); err != nil {
			log.Fatalf("Failed to subscribe to replication events: %v", err)
		}

		logger.Info("replication_enabled", "region", cfg.Region, "server_id", cfg.ServerID)
	}

	// Initialize supporting services
	presenceManager := presence.NewManager()
	pubsubService := pubsub.NewMemoryPubSub()
	dispatcher := webhook.NewDispatcher(cfg.WebhookURL)

	// Initialize Handlers
	srv := server.NewHandler(crdtEngine, presenceManager, pubsubService, dispatcher)
	healthChecker := server.NewHealthChecker(stateStore)

	// Router
	mux := http.NewServeMux()

	// Health Check Routes (no auth required)
	// Used by K8s liveness/readiness probes
	mux.HandleFunc("/healthz", healthChecker.HandleHealthz)
	mux.HandleFunc("/readyz", healthChecker.HandleReadyz)

	// API Routes
	mux.HandleFunc("/v1/sync/", srv.HandleWebSocket)
	mux.HandleFunc("/v1/presence/", srv.HandleGetPresence)
	mux.HandleFunc("/v1/stats", srv.HandleGetStats)
	mux.HandleFunc("/v1/history/", srv.HandleGetHistory)

	// Metrics Endpoint (P0 Enterprise Feature)
	mux.Handle("/metrics", promhttp.Handler())

	// Apply Middleware
	handler := auth.Middleware(mux)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: handler,
	}

	// Channel for server errors
	serverErrors := make(chan error, 1)

	// Start server
	go func() {
		logger.Info("server_starting",
			"port", cfg.Port,
			"strategy", strategy.Name(),
		)
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Wait for shutdown signal
	// We MUST capture SIGTERM to allow connections to drain.
	// K8s sends SIGTERM -> wait -> SIGKILL.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		logger.Info("shutdown_initiated", "signal", sig.String())

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := httpServer.Shutdown(ctx); err != nil {
			// CRITICAL: If this happens, some data in memory might not be flushed to BadgerDB.
			logger.Error("graceful_shutdown_failed", "error", err)
			httpServer.Close()
		}

		logger.Info("server_shutdown_complete")
	}
}
