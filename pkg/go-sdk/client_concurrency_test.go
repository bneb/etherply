package etherply_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	etherply "github.com/bneb/etherply/pkg/go-sdk"
	"github.com/gorilla/websocket"
)

// TestClient_Concurrency verifies that the Client methods are thread-safe.
// It runs with -race to detect data races.
func TestClient_Concurrency(t *testing.T) {
	// 1. Setup Mock Server
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			// Echo back
			err = c.WriteMessage(mt, message)
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	client := etherply.NewClient(wsURL, "mock-token")

	// Create a wait group for coordination
	var wg sync.WaitGroup
	start := make(chan struct{})

	// Action 1: Connect and Reconnect
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		// Hammer Connect
		for i := 0; i < 10; i++ {
			_ = client.Connect("ws-1")
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Action 2: Send Operations
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		for i := 0; i < 100; i++ {
			_ = client.SendOperation("key", fmt.Sprintf("val-%d", i))
			// Small sleep to allow interplay
			time.Sleep(2 * time.Millisecond)
		}
	}()

	// Action 3: Listen (Reads)
	// We need to call Listen once ideally, but concurrent calls shouldn't panic (though weird usage)
	// The typical usage is one Listen loop.
	// But let's verify that Listen doesn't race with Connect/Close.
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		// We just start listener and let it run
		client.Listen(func(msg map[string]interface{}) {
			// No-op handler
		})
	}()

	// Action 4: Close
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		time.Sleep(50 * time.Millisecond) // Let some ops pass
		_ = client.Close()
		time.Sleep(50 * time.Millisecond)
		// Reconnect
		_ = client.Connect("ws-1")
	}()

	// Start chaos
	close(start)
	wg.Wait()

	// Cleanup
	client.Close()
}
