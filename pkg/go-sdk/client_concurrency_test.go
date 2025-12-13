package etherply_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	etherply "github.com/bneb/etherply/pkg/go-sdk"
	"github.com/gorilla/websocket"
)

// TestClient_Concurrency verifies that Client.SendOperation is thread-safe.
// It spawns multiple goroutines sending operations simultaneously.
// Run with 'go test -race' to detect data races.
func TestClient_Concurrency(t *testing.T) {
	// 1. Setup Mock WebSocket Server
	upgrader := websocket.Upgrader{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return 
		}
		defer c.Close()
		for {
			_, _, err := c.ReadMessage()
			if err != nil {
				break
			}
		}
	}))
	defer server.Close()

	// Convert http URL to ws URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// 2. Create Client and Connect
	client := etherply.NewClient(wsURL, "test-token")
	err := client.Connect("test-workspace")
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer client.Close()

	// 3. Concurrent Load Test
	var wg sync.WaitGroup
	params := struct {
		Workers int
		Ops     int
	}{
		Workers: 20,
		Ops:     100,
	}

	wg.Add(params.Workers)
	for i := 0; i < params.Workers; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < params.Ops; j++ {
				// We don't check error here as we might close connection early or server might be slow,
				// the main point is to check for RACES in client structure access.
				_ = client.SendOperation("key", "val")
			}
		}(i)
	}

	wg.Wait()
}
