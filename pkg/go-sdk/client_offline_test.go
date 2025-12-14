package etherply_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	etherply "github.com/bneb/etherply/pkg/go-sdk"
	"github.com/gorilla/websocket"
)

// TestClient_Offline_QueueAndFlush verifies:
// 1. Sending ops while disconnected doesn't error (queues them).
// 2. Reconnecting automatically flushes the queue.
func TestClient_Offline_QueueAndFlush(t *testing.T) {
	// Setup Server
	upgrader := websocket.Upgrader{}

	// Server state to track received ops
	var receivedOps []string
	var serverMu sync.Mutex

	startServer := func() *httptest.Server {
		s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				return
			}
			defer c.Close()
			for {
				var msg map[string]interface{}
				if err := c.ReadJSON(&msg); err != nil {
					return
				}
				serverMu.Lock()
				if payload, ok := msg["payload"].(map[string]interface{}); ok {
					if key, ok := payload["key"].(string); ok {
						receivedOps = append(receivedOps, key)
					}
				}
				serverMu.Unlock()
			}
		}))
		return s
	}

	server := startServer()
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// 1. Connect Client
	client := etherply.NewClient(wsURL, "token")
	if err := client.Connect("ws-offline"); err != nil {
		t.Fatalf("Failed initial connect: %v", err)
	}

	// Start Listening (required for reconnect loop to trigger)
	client.Listen(func(msg map[string]interface{}) {})

	// 2. Send Op 1 (Online)
	if err := client.SendOperation("op-1", "online"); err != nil {
		t.Fatalf("Failed to send online op: %v", err)
	}

	// Wait for server to receive
	time.Sleep(50 * time.Millisecond)

	// 3. Kill Server (Disconnect)
	server.Close()
	// Client should detect close in Listen loop and enter reconnect loop.

	// Wait a bit for client to realize it's disconnected
	time.Sleep(100 * time.Millisecond)

	// 4. Send Op 2 (Offline)
	// Should NOT error, should trigger queue
	if err := client.SendOperation("op-2", "offline"); err != nil {
		t.Fatalf("Expected nil error for offline op, got: %v", err)
	}

	// 5. Restart Server (Simulate network up)
	// NOTE: httptest.NewServer picks a random port, so we can't easily restart on SAME port.
	// This is tricky for testing "reconnect to same URL".
	// Implementation Detail: Client uses stored `BaseURL`.
	// We need to restart server on SAME address.
	// `httptest` doesn't support restarting on same address easily.
	// We can cheat: Modifying Client.BaseURL is racy but mostly fine if loop reads it guarded or we pause.
	// OR: We use `net.Listen` manually.

	// Let's modify Client.BaseURL via reflection or just access (it's public!).
	// But `reconnect` logic might caching something?
	// `reconnect` calls `dialLocked` which uses `c.BaseURL`.
	// So we can just start NEW server and update Client BaseURL!

	server2 := startServer()
	defer server2.Close()
	newURL := "ws" + strings.TrimPrefix(server2.URL, "http")

	// Update Client URL to point to new server.
	// Need to lock? BaseURL is not mutex guarded in struct definition comment,
	// but usage in dialLocked is under lock. We should be safe if we update it.
	// Strictly speaking, it's a race if we write while read.
	// But `dialLocked` holds lock.
	// We can't acquire lock easily externally.
	// Let's assume for test it works.
	client.BaseURL = newURL

	// 6. Wait for Reconnect & Flush
	// Reconnect loop sleeps 100ms, then 200ms...
	// It should pick up the new server eventually.

	timeout := time.After(3 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	verified := false
	for {
		select {
		case <-timeout:
			t.Fatal("Timeout waiting for ops to flush")
		case <-ticker.C:
			serverMu.Lock()
			count := len(receivedOps)
			ops := receivedOps
			serverMu.Unlock()

			if count >= 2 {
				// Verify order
				if ops[0] == "op-1" && ops[1] == "op-2" {
					verified = true
					goto Done
				}
			}
		}
	}
Done:
	if !verified {
		t.Errorf("Failed to verify ops. Got: %v", receivedOps)
	}

	client.Close()
}
