package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleGetStats(t *testing.T) {
	handler := createTestHandler()

	req := httptest.NewRequest("GET", "/v1/stats", nil)
	rr := httptest.NewRecorder()

	handler.HandleGetStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var stats map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatal(err)
	}

	// Verify Structure
	if _, ok := stats["connections"]; !ok {
		t.Error("Missing 'connections' field")
	}
	if _, ok := stats["persistence"]; !ok {
		t.Error("Missing 'persistence' field")
	}
	if _, ok := stats["server_time"]; !ok {
		t.Error("Missing 'server_time' field")
	}
}
