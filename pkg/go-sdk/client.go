package etherply

import (
	"log"

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

func (c *Client) SendOperation(key string, value interface{}) error {
	msg := map[string]interface{}{
		"type": "op",
		"payload": map[string]interface{}{
			"key":   key,
			"value": value,
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
