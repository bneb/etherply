package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/bneb/etherply/etherply-sync-server/internal/crdt"
	"github.com/bneb/etherply/etherply-sync-server/internal/presence"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	StartBufferSize: 4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for Demo/DX
	},
}

type Handler struct {
	crdtEngine      *crdt.Engine
	presenceManager *presence.Manager
	// Simple Hub for broadcasting (basic implementation)
	hubMu sync.Mutex
	hub   map[string]map[*websocket.Conn]bool // workspaceID -> conns
}

func NewHandler(e *crdt.Engine, p *presence.Manager) *Handler {
	return &Handler{
		crdtEngine:      e,
		presenceManager: p,
		hub:             make(map[string]map[*websocket.Conn]bool),
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
		log.Printf("[WS] Upgrade error: %v", err)
		return
	}

	// Register Connection
	h.registerConnection(workspaceID, conn)
	h.presenceManager.AddUser(workspaceID, presence.User{UserID: userID, Status: "online"})
	
	log.Printf("[METRIC] concurrent_connection_peak | workspace=%s user_tier=FREE", workspaceID)

	defer func() {
		h.unregisterConnection(workspaceID, conn)
		h.presenceManager.RemoveUser(workspaceID, userID)
		conn.Close()
	}()

	// Send Initial State
	state, err := h.crdtEngine.GetFullState(workspaceID)
	if err == nil {
		// Wrap in a sync message
		msg := map[string]interface{}{
			"type": "init",
			"data": state,
		}
		conn.WriteJSON(msg)
	}

	for {
		// Read Message
		var rawMsg map[string]interface{}
		if err := conn.ReadJSON(&rawMsg); err != nil {
			break
		}

		// Handle "ping" or "op"
		msgType, _ := rawMsg["type"].(string)
		
		if msgType == "op" {
			// Parse Operation
			// Expecting payload structure matching crdt.Operation
			// We cheat a bit with JSON marshalling here for the MVP loop
			payloadBytes, _ := json.Marshal(rawMsg["payload"])
			var op crdt.Operation
			json.Unmarshal(payloadBytes, &op)
			
			// Process
			op.WorkspaceID = workspaceID // Force security
			err := h.crdtEngine.ProcessOperation(op)
			if err != nil {
				log.Printf("Error processing op: %v", err)
				continue
			}

			// Broadcast to others in workspace
			h.broadcast(workspaceID, rawMsg, conn)
		}
	}
}

func (h *Handler) registerConnection(workspaceID string, conn *websocket.Conn) {
	h.hubMu.Lock()
	defer h.hubMu.Unlock()
	if _, ok := h.hub[workspaceID]; !ok {
		h.hub[workspaceID] = make(map[*websocket.Conn]bool)
	}
	h.hub[workspaceID][conn] = true
}

func (h *Handler) unregisterConnection(workspaceID string, conn *websocket.Conn) {
	h.hubMu.Lock()
	defer h.hubMu.Unlock()
	if _, ok := h.hub[workspaceID]; ok {
		delete(h.hub[workspaceID], conn)
		if len(h.hub[workspaceID]) == 0 {
			delete(h.hub, workspaceID)
		}
	}
}

func (h *Handler) broadcast(workspaceID string, msg interface{}, sender *websocket.Conn) {
	h.hubMu.Lock()
	defer h.hubMu.Unlock()
	
	if clients, ok := h.hub[workspaceID]; ok {
		for client := range clients {
			if client != sender { // Don't echo back if client handles optimistic UI (usually)
				// For simple CRDTs we might echo back to confirm, but typically we don't for bandwidth.
				err := client.WriteJSON(msg)
				if err != nil {
					log.Printf("Broadcast error: %v", err)
					client.Close()
					delete(clients, client)
				}
			}
		}
	}
}
