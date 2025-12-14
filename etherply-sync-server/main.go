package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/bneb/etherply/etherply-sync-server/internal/pubsub"
	"github.com/bneb/etherply/etherply-sync-server/internal/replication"
	"github.com/bneb/etherply/etherply-sync-server/internal/server"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
	"github.com/bneb/etherply/etherply-sync-server/internal/webhook"
)

// main is the entry point for the EtherPly Sync Server.
//
// Architecture Overview:
//
//  1. Config: Reads env vars (PORT, JWT_SECRET, SHUTDOWN_TIMEOUT_SECONDS).
//  2. Persistence: Initializes BadgerDB for durability.
//  3. Engine: Starts the CRDT Engine which manages state per workspace.
//  4. Replication: Optionally enables multi-region replication via NATS.
//  5. Presence: Starts the ephemeral Presence Manager.
//  6. HTTP: Sets up routes and health checks.
//  7. Graceful Shutdown: Handles SIGTERM/SIGINT for clean connection draining.
//
// This server is designed to be stateless regarding "sessions" but stateful regarding "document data".
// It can be deployed to PaaS like Fly.io or Kubernetes with proper health probes.
func main() {
	// PORT environment variable is standard for Fly.io, Heroku, and Kubernetes.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Graceful shutdown timeout (default 30 seconds)
	shutdownTimeout := 30 * time.Second
	if timeoutStr := os.Getenv("SHUTDOWN_TIMEOUT_SECONDS"); timeoutStr != "" {
		if t, err := strconv.Atoi(timeoutStr); err == nil && t > 0 {
			shutdownTimeout = time.Duration(t) * time.Second
		}
	}

	// AUTH SECURITY
	// We aggressively fail if no secret is provided to prevent accidental insecure deployments.
	jwtSecret := os.Getenv("ETHERPLY_JWT_SECRET")
	if jwtSecret == "" {
		log.Println("[CRITICAL] ETHERPLY_JWT_SECRET is not set.")
		log.Fatal("Server cannot start in unsecure mode. Please set ETHERPLY_JWT_SECRET.")
	}
	auth.Init(jwtSecret)

	// Initialize Store (BadgerDB v4 for Production-Grade Persistence)
	// BADGER_PATH allows customization for container deployments
	badgerPath := os.Getenv("BADGER_PATH")
	if badgerPath == "" {
		badgerPath = "./badger.db"
	}
	stateStore, err := store.NewBadgerStore(badgerPath)
	if err != nil {
		log.Fatalf("Failed to initialize persistence layer at %s: %v", badgerPath, err)
	}
	defer stateStore.Close()
	log.Printf("Persistence layer initialized at: %s", badgerPath)

	// Initialize CRDT Engine
	crdtEngine := crdt.NewEngine(stateStore)

	// Initialize Multi-Region Replication (if configured)
	// NATS_URL: comma-separated list of NATS server URLs
	// REGION: geographic region identifier (e.g., "us-east-1")
	// SERVER_ID: unique identifier for this server instance
	var replicator *replication.NATSReplicator
	natsURL := os.Getenv("NATS_URL")
	if natsURL != "" {
		region := os.Getenv("REGION")
		if region == "" {
			region = "default"
		}
		serverID := os.Getenv("SERVER_ID")
		if serverID == "" {
			serverID = "sync-server-" + port
		}

		natsURLs := strings.Split(natsURL, ",")
		replicator, err = replication.NewNATSReplicator(replication.Config{
			ServerID: serverID,
			Region:   region,
			NATSURLs: natsURLs,
		})
		if err != nil {
			log.Fatalf("Failed to initialize NATS replicator: %v", err)
		}
		defer replicator.Close()

		// Wire replication to CRDT engine
		crdtEngine.SetReplicator(replicator, region, serverID)

		// Subscribe to incoming changes from peer replicas
		if err := replicator.Subscribe(func(event replication.ChangeEvent) error {
			return crdtEngine.ApplyRemoteChanges(event.WorkspaceID, event.Changes)
		}); err != nil {
			log.Fatalf("Failed to subscribe to replication events: %v", err)
		}

		log.Printf("Multi-region replication enabled: region=%s, server_id=%s", region, serverID)
	}

	// Initialize Presence Manager (Ephemeral, In-Memory)
	presenceManager := presence.NewManager()

	// Initialize Pub/Sub Layer (In-Memory for now, replaceable with Redis later)
	pubsubService := pubsub.NewMemoryPubSub()

	// Initialize Webhook Dispatcher
	webhookURL := os.Getenv("WEBHOOK_URL")
	dispatcher := webhook.NewDispatcher(webhookURL)

	// Initialize Handlers
	srv := server.NewHandler(crdtEngine, presenceManager, pubsubService, dispatcher)
	healthChecker := server.NewHealthChecker(stateStore)

	// Router
	mux := http.NewServeMux()

	// Health Check Routes (no auth required - Kubernetes probes)
	mux.HandleFunc("/healthz", healthChecker.HandleHealthz)
	mux.HandleFunc("/readyz", healthChecker.HandleReadyz)

	// API Routes
	mux.HandleFunc("/v1/sync/", srv.HandleWebSocket)       // WS: /v1/sync/{workspace_id}
	mux.HandleFunc("/v1/presence/", srv.HandleGetPresence) // REST: /v1/presence/{workspace_id}
	mux.HandleFunc("/v1/stats", srv.HandleGetStats)        // REST: /v1/stats
	mux.HandleFunc("/v1/history/", srv.HandleGetHistory)   // REST: /v1/history/{workspace_id}

	// Apply Middleware (Logging, Auth, Recovery)
	handler := auth.Middleware(mux)

	// Create HTTP server with graceful shutdown support
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	// Channel to signal server errors
	serverErrors := make(chan error, 1)

	// Start server in goroutine
	go func() {
		log.Printf("EtherPly Sync Server starting on port %s", port)
		log.Printf("Health endpoints: /healthz (liveness), /readyz (readiness)")
		serverErrors <- httpServer.ListenAndServe()
	}()

	// Wait for interrupt signal or server error
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("Received signal %v, initiating graceful shutdown...", sig)

		// Create shutdown context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v, forcing close", err)
			httpServer.Close()
		}

		log.Println("Server shutdown complete")
	}
}
