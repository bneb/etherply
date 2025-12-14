package etherply

import (
	"fmt"
	"log"
	"sync"
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
	mu      sync.RWMutex // Guards writes to Conn and Queue

	// Offline Support
	queue       []map[string]interface{}
	workspaceID string
	isClosed    bool // True if user explicitly called Close()
}

// NewClient creates a new Client instance.
// baseURL should be the root URL of the sync server (e.g., "ws://localhost:8080").
// token must be a valid JWT signed with the server's secret.
func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		queue:   make([]map[string]interface{}, 0),
	}
}

// Connect establishes the WebSocket connection to a specific workspace.
// It performs the handshake and authentication.
// Returns an error if the server is unreachable or the token is invalid (401).
func (c *Client) Connect(workspaceID string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.workspaceID = workspaceID
	c.isClosed = false

	if err := c.dialLocked(); err != nil {
		return err
	}

	log.Println("[SDK] Connected to EtherPly Sync Server")
	return nil
}

// internal helper, assumes lock is held
func (c *Client) dialLocked() error {
	// Construct URL
	// Query params are used for auth during WS handshake.
	url := c.BaseURL + "/v1/sync/" + c.workspaceID + "?token=" + c.Token

	// DefaultDialer is used; in production, you might want to customize the HandshakeTimeout.
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return fmt.Errorf("failed to dial %s: %w", url, err)
	}
	c.Conn = conn
	return nil
}

// SendOperation transmits a key-value operation to the sync server.
// It automatically attaches a microsecond-precision timestamp for LWW (Last-Write-Wins)
// conflict resolution.
//
// Offline Support: If the client is disconnected, the operation is queued and sent
// automatically upon reconnection.
func (c *Client) SendOperation(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	msg := map[string]interface{}{
		"type": "op",
		"payload": map[string]interface{}{
			"key":       key,
			"value":     value,
			"timestamp": time.Now().UnixMicro(),
		},
	}

	// If disconnected or closed (but queueing allowed?), just queue if generic disconnect.
	// We want to queue if NOT closed by user, but connection is dropped.
	// But Connect() resets isClosed to false.
	if c.Conn == nil {
		if c.isClosed {
			return fmt.Errorf("client is closed")
		}
		// Offline buffering
		c.queue = append(c.queue, msg)
		return nil
	}

	err := c.Conn.WriteJSON(msg)
	if err != nil {
		// Network error during write? Queue it!
		log.Printf("[SDK] Write failed, queuing operation: %v", err)
		c.queue = append(c.queue, msg)
		// We could mark c.Conn = nil here to force reconnect logic elsewhere
		c.Conn.Close()
		c.Conn = nil
		return nil // Optimistic success
	}

	return nil
}

// Listen starts a background goroutine to read messages from the server.
// The handler function is called for every message received.
//
// Blocking: This function itself returns immediately, but the spawned goroutine runs until
// the connection is closed.
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
			c.mu.RLock()
			conn := c.Conn
			closed := c.isClosed
			c.mu.RUnlock()

			// If explicitly closed, stop.
			if closed {
				return
			}

			// If disconnected, try to reconnect
			if conn == nil {
				if err := c.reconnect(); err != nil {
					// Reconnect failed (still failing), backoff is handled inside reconnect.
					// If we are here, reconnect returned error which might mean fatal or just "try again loop".
					// We'll just continue loop to try again.
					time.Sleep(1 * time.Second) // Base sleep if reconnect returns error instantly
					continue
				}
				// Reconnected! Loop again to pick up new conn
				continue
			}

			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				// Normal closure or error
				log.Printf("[SDK] Connection read error: %v. Reconnecting...", err)

				// Mark as disconnected so next loop triggers reconnect
				c.mu.Lock()
				if c.Conn == conn { // Only nullify if it hasn't changed (race check)
					c.Conn = nil
					// conn is the local value, safe to close it (it might be c.Conn before nil-ing)
					// But we should close the connection we just failed on.
					// Actually c.Conn was nilled, so we can't call Close on it.
					// We should Close the `conn` variable we hold.
				}
				c.mu.Unlock()
				conn.Close() // Ensure closed
				continue
			}
			handler(msg)
		}
	}()
}

// reconnect attempts to re-establish connection with exponential backoff.
// It also flushes the operation queue upon success.
func (c *Client) reconnect() error {
	backoff := 100 * time.Millisecond
	maxBackoff := 10 * time.Second

	for {
		// Check if user closed while we were waiting
		c.mu.RLock()
		if c.isClosed {
			c.mu.RUnlock()
			return nil
		}
		c.mu.RUnlock()

		log.Println("[SDK] Attempting to reconnect...")

		c.mu.Lock()
		err := c.dialLocked()
		if err == nil {
			// Success! Flush queue.
			log.Printf("[SDK] Reconnected. Flushing %d operations...", len(c.queue))

			// We iterate queue and try to send.
			// Check: if send fails again during flushing?
			// We should keep them in queue.
			failedCount := 0
			newQueue := make([]map[string]interface{}, 0)

			for _, atomicOp := range c.queue {
				// We write directly to Conn because SendOperation would just re-queue on nil
				// But we HAVE conn here (locked).
				wErr := c.Conn.WriteJSON(atomicOp)
				if wErr != nil {
					log.Printf("[SDK] Failed to flush op: %v", wErr)
					newQueue = append(newQueue, atomicOp)
					failedCount++
				}
			}
			c.queue = newQueue

			c.mu.Unlock()

			if failedCount > 0 {
				log.Printf("[SDK] Failed to flush %d ops. Will retry next cycle.", failedCount)
				// If we failed to flush, maybe connection is bad.
				// The listener loop will pick this up when it tries to read and fails, or simply we hold open?
				// For now, assume good enough.
			}
			return nil
		}
		c.mu.Unlock()

		// Failure
		log.Printf("[SDK] Reconnect failed: %v. Retrying in %v", err, backoff)
		time.Sleep(backoff)

		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

// Close terminates the WebSocket connection gracefully.
// It sends a close message to the server and closes the underlying connection.
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.isClosed = true // Signal to stop reconnecting

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
