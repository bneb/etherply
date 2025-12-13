package main

import (
	"log"
	"net/http"
	"os"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/bneb/etherply/etherply-sync-server/internal/server"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

// main is the entry point for the EtherPly Sync Server.
//
// Architecture Overview:
//  1. Config: Reads env vars (PORT, JWT_SECRET).
//  2. Persistence: Initializes DiskStore (AOF) for durability.
//  3. Engine: Starts the CRDT Engine which manages state per workspace.
//  4. Presence: Starts the ephemeral Presence Manager.
//  5. HTTP: Sets up routes and blocked/non-blocking handlers.
//
// This server is designed to be stateless regarding "sessions" but stateful regarding "document data".
// It can be deployed to PaaS like Fly.io or Heroku, but requires a persistent volume for the AOF file
// if data durability across restarts is required.
func main() {
	// PORT environment variable is standard for Fly.io and Heroku.
	// If missing, we default to 8080 for local development convenience.
	// 2 AM Rule: Do not assume standard ports are free. Check lsof -i :8080 if failing.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
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
	// We use a local directory "badger.db". In production, this path comes from env.
	// Ideally this should be on a mounted volume.
	stateStore, err := store.NewBadgerStore("badger.db")
	if err != nil {
		log.Fatalf("Failed to initialize persistence layer: %v", err)
	}
	defer stateStore.Close()

	// Initialize CRDT Engine
	// The Engine holds the "Truth" of the document state in memory.
	crdtEngine := crdt.NewEngine(stateStore)

	// Initialize Presence Manager (Ephemeral, In-Memory)
	presenceManager := presence.NewManager()

	// Initialize Server Handler
	srv := server.NewHandler(crdtEngine, presenceManager)

	// Router
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/v1/sync/", srv.HandleWebSocket)       // WS: /v1/sync/{workspace_id}
	mux.HandleFunc("/v1/presence/", srv.HandleGetPresence) // REST: /v1/presence/{workspace_id}

	// Apply Middleware (Logging, Auth, Recovery)
	handler := auth.Middleware(mux)

	log.Printf("EtherPly Sync Server starting on port %s", port)
	// ListenAndServe is blocking.
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
