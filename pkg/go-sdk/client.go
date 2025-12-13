package etherply

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	BaseURL string
	Token   string
	Conn    *websocket.Conn
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
	}
}

func (c *Client) Connect(workspaceID string) error {
	// Construct URL
	url := c.BaseURL + "/v1/sync/" + workspaceID + "?token=" + c.Token
	
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.Conn = conn
	log.Println("[SDK] Connected to EtherPly Sync Server")
	return nil
}

// SendOperation transmits a key-value operation to the sync server.
// It includes a timestamp for LWW (Last-Write-Wins) conflict resolution.
// Returns an error if the connection has not been established via Connect().
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
	return c.Conn.WriteJSON(msg)
}

func (c *Client) Listen(handler func(msg map[string]interface{})) {
	go func() {
		for {
			var msg map[string]interface{}
			err := c.Conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("[SDK] Read error: %v", err)
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
	// Send close message per WebSocket protocol
	err := c.Conn.WriteMessage(websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Printf("[SDK] Error sending close message: %v", err)
	}
	return c.Conn.Close()
}
