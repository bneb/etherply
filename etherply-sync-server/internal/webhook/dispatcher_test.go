package webhook_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/webhook"
)

func TestDispatcher_Dispatch(t *testing.T) {
	var mu sync.Mutex
	var received []webhook.EventPayload

	// 1. Start Mock Server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}

		var payload webhook.EventPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Errorf("Failed to decode body: %v", err)
		}

		mu.Lock()
		received = append(received, payload)
		mu.Unlock()

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// 2. Init Dispatcher
	d := webhook.NewDispatcher(ts.URL)

	// 3. Dispatch Events (Async)
	d.Dispatch("test.event", map[string]string{"foo": "bar"})
	d.Dispatch("test.event.2", "simple string")

	// 4. Wait for processing (worker is async)
	// Simple polling
	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		mu.Lock()
		count := len(received)
		mu.Unlock()
		if count >= 2 {
			break
		}
		time.Sleep(10 * time.Millisecond)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(received) != 2 {
		t.Fatalf("Expected 2 events, got %d", len(received))
	}

	if received[0].Event != "test.event" {
		t.Errorf("Expected test.event, got %s", received[0].Event)
	}

	// Verify timestamp exists
	if received[0].Timestamp == 0 {
		t.Error("Timestamp should be set")
	}
}
