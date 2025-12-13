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

func main() {
	// PORT environment variable is standard for Fly.io and Heroku.
	// If missing, we default to 8080 for local development convenience.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize Store (Persistent DiskStore for MVP Durability)
	// We use a local file "etherply.aof". In production, this path comes from env.
	stateStore, err := store.NewDiskStore("etherply.aof")
	if err != nil {
		log.Fatalf("Failed to initialize persistence layer: %v", err)
	}
	defer stateStore.Close()

	// Initialize CRDT Engine
	crdtEngine := crdt.NewEngine(stateStore)

	// Initialize Presence Manager
	presenceManager := presence.NewManager()

	// Initialize Server Handler
	srv := server.NewHandler(crdtEngine, presenceManager)

	// Router
	mux := http.NewServeMux()

	// Public Routes
	mux.HandleFunc("/v1/sync/", srv.HandleWebSocket)       // /v1/sync/{workspace_id}
	mux.HandleFunc("/v1/presence/", srv.HandleGetPresence) // /v1/presence/{workspace_id}

	// Apply Middleware
	handler := auth.Middleware(mux)

	log.Printf("EtherPly Sync Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
