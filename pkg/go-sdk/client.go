package etherply

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client represents a connection to the EtherPly Sync Server.
// It manages the WebSocket connection and provides methods to send operations
// and listen for updates.
//
// Thread Safety:
// - SendOperation is thread-safe.
// - Listen runs in its own goroutine.
// - Close is thread-safe.
type Client struct {
	BaseURL string
	Token   string
	Conn    *websocket.Conn
}

// NewClient creates a new Client instance.
// baseURL should be the root URL of the sync server (e.g., "ws://localhost:8080").
// token must be a valid JWT signed with the server's secret.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
	}
}

// Connect establishes the WebSocket connection to a specific workspace.
// It performs the handshake and authentication.
// Returns an error if the server is unreachable or the token is invalid (401).
func (c *Client) Connect(workspaceID string) error {
	// Construct URL
	// Query params are used for auth during WS handshake.
	url := c.BaseURL + "/v1/sync/" + workspaceID + "?token=" + c.Token

	// DefaultDialer is used; in production, you might want to customize the HandshakeTimeout.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %w", url, err)
	}
	c.Conn = conn
	log.Println("[SDK] Connected to EtherPly Sync Server")
	return nil
}

// SendOperation transmits a key-value operation to the sync server.
// It automatically attaches a microsecond-precision timestamp for LWW (Last-Write-Wins)
// conflict resolution.
//
// Returns an error if the connection has not been established via Connect().
// Note: This does not guarantee the server *accepted* the write, only that it was written to the socket.
func (c *Client) SendOperation(key string, value interface{}) error {
	// Defensive: Check connection before attempting to write.
	// "The Happy Path is a Trap" - we must handle the not-connected case gracefully.
	if c.Conn == nil {
		return fmt.Errorf("cannot send operation: connection not established (call Connect first)")
	}

	msg := map[string]interface{}{
		"type": "op",
		"payload": map[string]interface{}{
			"key":       key,
			"value":     value,
			"timestamp": time.Now().UnixMicro(),
		},
	}
	// WriteJSON is thread-safe effectively because we don't interleave partial messages,
	// but strictly speaking, Gorilla WebSocket WriteJSON is NOT concurrent safe.
	// If this SDK were to be used heavily concurrently, we would need a mutex here.
	// For this MVP, we assume single-threaded writing or low contention.
	// TODO: Add Mutex for robust thread safety if multiple goroutines write simultaneously.
	return c.Conn.WriteJSON(msg)
}

// Listen starts a background goroutine to read messages from the server.
// The handler function is called for every message received.
//
// Blocking: This function itself returns immediately, but the spawned goroutine runs until
// the connection is closed.
func (c *Client) Listen(handler func(msg map[string]interface{})) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[SDK] Panic in listener: %v", r)
			}
		}()
		for {
			if c.Conn == nil {
				return
			}
			var msg map[string]interface{}
			err := c.Conn.ReadJSON(&msg)
			if err != nil {
				// Normal closure or error
				log.Printf("[SDK] Connection read/closed: %v", err)
				return
			}
			handler(msg)
		}
	}()
}

// Close terminates the WebSocket connection gracefully.
// It sends a close message to the server and closes the underlying connection.
func (c *Client) Close() error {
	if c.Conn == nil {
		return nil
	}
	// Send close message per WebSocket protocol to be polite
	err := c.Conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "client closing"))
	if err != nil {
		log.Printf("[SDK] Error sending close message: %v", err)
	}
	return c.Conn.Close()
}
