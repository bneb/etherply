// Package server provides HTTP and WebSocket handlers for the EtherPly sync engine.
// It orchestrates connections between clients and the CRDT engine, handling:
//   - WebSocket upgrade and lifecycle management
//   - Message routing (operations broadcast to workspace members)
//   - Presence tracking integration
//   - Initial state synchronization for new connections
package server

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/bneb/etherply/etherply-sync-server/internal/pubsub"
	"github.com/bneb/etherply/etherply-sync-server/internal/webhook"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for Demo/DX
	},
}

type Handler struct {
	crdtEngine      *crdt.Engine
	presenceManager *presence.Manager
	pubsub          pubsub.PubSub
	webhook         *webhook.Dispatcher
	logger          *slog.Logger
}

func NewHandler(e *crdt.Engine, p *presence.Manager, ps pubsub.PubSub, wh *webhook.Dispatcher) *Handler {
	// Default to JSON handler for structured output, writing to stderr
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	return &Handler{
		crdtEngine:      e,
		presenceManager: p,
		pubsub:          ps,
		webhook:         wh,
		logger:          logger,
	}
}

func (h *Handler) HandleGetPresence(w http.ResponseWriter, r *http.Request) {
	// Path: /v1/presence/{workspace_id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	workspaceID := parts[3]

	users := h.presenceManager.GetUsers(workspaceID)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *Handler) HandleGetStats(w http.ResponseWriter, r *http.Request) {
	pubsubStats := h.pubsub.Stats()
	engineStats, err := h.crdtEngine.Stats()

	if err != nil {
		h.logger.Error("stats_failed", slog.Any("error", err))
		http.Error(w, "Failed to retrieve stats", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"connections": pubsubStats,
		"persistence": engineStats,
		"server_time": time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	// Path: /v1/history/{workspace_id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	workspaceID := parts[3]

	// Auth check implicitly done by middleware, but scoped check?
	// If strictly enforcing scopes, we check for "read".
	// But middleware passes if token valid.
	// For "Ironclad" security, we could check scopes here too.
	// Legacy/Dev: Allow.

	history, err := h.crdtEngine.GetHistory(workspaceID)
	if err != nil {
		h.logger.Error("history_failed", slog.Any("error", err))
		http.Error(w, "Failed to retrieve history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func (h *Handler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Path: /v1/sync/{workspace_id}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	workspaceID := parts[3]

	// Stub User ID (normally from Auth)
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		userID = "anon"
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("ws_upgrade_failed", slog.Any("error", err))
		return
	}

	// 1. Subscribe to PubSub
	rxChan, unsub := h.pubsub.Subscribe(workspaceID)

	// Register Connection Presence
	h.presenceManager.AddUser(workspaceID, presence.User{UserID: userID, Status: presence.StatusOnline})

	// Webhook: client.connected
	h.webhook.Dispatch("client.connected", map[string]string{
		"workspace_id": workspaceID,
		"user_id":      userID,
	})

	h.logger.Info("concurrent_connection_peak",
		slog.String("event", "metric"),
		slog.String("workspace_id", workspaceID),
		slog.String("user_tier", "FREE"),
	)

	// Clean up on exit
	defer func() {
		unsub()
		h.presenceManager.RemoveUser(workspaceID, userID)
		// Webhook: client.disconnected
		h.webhook.Dispatch("client.disconnected", map[string]string{
			"workspace_id": workspaceID,
			"user_id":      userID,
		})
		conn.Close()
	}()

	// 2. Start Writer Goroutine (WritePump)
	// Consumes messages from PubSub and writes to WebSocket
	go func() {
		defer conn.Close() // Ensure close if this routine exits
		for msg := range rxChan {
			// Don't echo back if senderID matches?
			// The current impl of PubSub stores SenderID in Message.
			// But wait, the WebSocket itself doesn't have a unique ID unless we assign one.
			// `conn` pointer is unique address but not serializable easily.
			// Using userID? But users can have multiple tabs.
			// For now, simple echo or filter if we add ID.

			// If we want basic functionality: Just echo to all. Client can filter echoes if needed.
			// But usually we want to avoid echo.
			// Let's assume we send everything for now (simplest transition).

			// We receive raw bytes in Payload
			// But WebSocket expects JSON object usually if we rely on `conn.WriteJSON`.
			// Payload is []byte.
			// `conn.WriteMessage(websocket.TextMessage, msg.Payload)`

			// Wait, previous broadcast filtered sender.
			// If we send back to sender, they might apply double or ignore.
			// Automerge handles idempotency, so echoes are fine logically but wasteful bandwidth.

			err := conn.WriteMessage(websocket.TextMessage, msg.Payload)
			if err != nil {
				return // Stop writer if write fails
			}
		}
	}()

	// 3. Send Initial State
	snapshot, err := h.crdtEngine.GetFullState(workspaceID)
	if err == nil {
		// Wrap in a sync message
		msg := map[string]interface{}{
			"type":  "init",
			"data":  snapshot.Data,
			"heads": snapshot.Heads,
		}
		conn.WriteJSON(msg)
	}

	// 4. Read Loop (Main routine blocks here)
	for {
		// Read Message
		// rawMsg is map
		var rawMsg map[string]interface{}
		if err := conn.ReadJSON(&rawMsg); err != nil {
			break
		}

		// Handle "ping" or "op"
		msgType, _ := rawMsg["type"].(string)

		if msgType == "op" {
			// Parse Operation
			payloadBytes, _ := json.Marshal(rawMsg["payload"])
			var op crdt.Operation
			json.Unmarshal(payloadBytes, &op)

			// Process
			op.WorkspaceID = workspaceID // Force security

			// ACL Check: "write" scope
			// Legacy/Dev: If no scopes defined in token, allow all.
			// If scopes defined, must have "write".
			scopes := auth.ScopesFromContext(r.Context())
			if len(scopes) > 0 {
				canWrite := false
				for _, s := range scopes {
					if s == "write" || s == "admin" {
						canWrite = true
						break
					}
				}
				if !canWrite {
					h.logger.Warn("acl_denied", slog.String("reason", "missing_write_scope"), slog.String("user_id", userID))
					// Send error to client?
					conn.WriteJSON(map[string]interface{}{
						"type":    "error",
						"payload": "permission_denied: missing 'write' scope",
					})
					continue
				}
			}

			err := h.crdtEngine.ProcessOperation(op)
			if err != nil {
				h.logger.Error("op_processing_failed", slog.Any("error", err))
				continue
			}

			// Webhook: doc.updated
			h.webhook.Dispatch("doc.updated", map[string]string{
				"workspace_id": workspaceID,
				"user_id":      userID,
				"key":          op.Key,
			})

			// Broadcast via PubSub
			// We need to re-serialize the full message to send to others
			fullMsgBytes, _ := json.Marshal(rawMsg)

			h.pubsub.Publish(workspaceID, pubsub.Message{
				Topic:   workspaceID,
				Payload: fullMsgBytes,
				// SenderID: ???
			})
		}
	}
}
