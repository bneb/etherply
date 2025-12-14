package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bneb/etherply/etherply-sync-server/internal/server"
	"github.com/bneb/etherply/etherply-sync-server/internal/store"
)

func TestHealthz_ReturnsOK(t *testing.T) {
	memStore := store.NewMemoryStore()
	healthChecker := server.NewHealthChecker(memStore)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	healthChecker.HandleHealthz(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response server.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.Timestamp == "" {
		t.Error("Expected timestamp to be set")
	}

	if response.Uptime == "" {
		t.Error("Expected uptime to be set")
	}

	t.Logf("Health response: %+v", response)
}

func TestReadyz_HealthyStore_ReturnsOK(t *testing.T) {
	memStore := store.NewMemoryStore()
	healthChecker := server.NewHealthChecker(memStore)

	req := httptest.NewRequest("GET", "/readyz", nil)
	rr := httptest.NewRecorder()

	healthChecker.HandleReadyz(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	var response server.HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Status != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", response.Status)
	}

	if response.Checks["store"] != "ok" {
		t.Errorf("Expected store check 'ok', got '%s'", response.Checks["store"])
	}

	t.Logf("Ready response: %+v", response)
}

func TestReadyz_ContentType(t *testing.T) {
	memStore := store.NewMemoryStore()
	healthChecker := server.NewHealthChecker(memStore)

	req := httptest.NewRequest("GET", "/readyz", nil)
	rr := httptest.NewRecorder()

	healthChecker.HandleReadyz(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHealthz_ContentType(t *testing.T) {
	memStore := store.NewMemoryStore()
	healthChecker := server.NewHealthChecker(memStore)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rr := httptest.NewRecorder()

	healthChecker.HandleHealthz(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}
