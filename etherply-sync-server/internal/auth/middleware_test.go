package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bneb/etherply/etherply-sync-server/internal/auth"
	"github.com/golang-jwt/jwt/v5"
)

// generateTestToken creates a valid JWT for testing the middleware.
// Uses the same HS256 signing method expected by the auth package.
func generateTestToken(secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(time.Hour).Unix(),
	})
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestMiddleware_CORSPreflight(t *testing.T) {
	// Setup: OPTIONS requests should return 200 with CORS headers, no auth needed
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called for OPTIONS request")
	}))

	req := httptest.NewRequest("OPTIONS", "/v1/sync/test-workspace", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}

	// Verify CORS headers are set
	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Missing CORS Allow-Origin header")
	}
	if rr.Header().Get("Access-Control-Allow-Methods") == "" {
		t.Error("Missing CORS Allow-Methods header")
	}
	if rr.Header().Get("Access-Control-Allow-Headers") == "" {
		t.Error("Missing CORS Allow-Headers header")
	}
}

func TestMiddleware_MissingToken(t *testing.T) {
	auth.Init("test-secret")

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called without token")
	}))

	req := httptest.NewRequest("GET", "/v1/presence/test-workspace", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", rr.Code)
	}
}

func TestMiddleware_InvalidToken(t *testing.T) {
	auth.Init("test-secret")

	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called with invalid token")
	}))

	req := httptest.NewRequest("GET", "/v1/presence/test-workspace", nil)
	req.Header.Set("Authorization", "Bearer invalid-garbage-token")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status 401 Unauthorized, got %d", rr.Code)
	}
}

func TestMiddleware_ValidToken_Header(t *testing.T) {
	secret := "test-secret-for-valid-header"
	auth.Init(secret)

	validToken := generateTestToken(secret)

	handlerCalled := false
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/v1/presence/test-workspace", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("Handler should have been called with valid token in header")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}

func TestMiddleware_ValidToken_QueryParam(t *testing.T) {
	secret := "test-secret-for-query-param"
	auth.Init(secret)

	validToken := generateTestToken(secret)

	handlerCalled := false
	handler := auth.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		w.WriteHeader(http.StatusOK)
	}))

	// Token passed via query parameter (common for WebSocket connections)
	req := httptest.NewRequest("GET", "/v1/sync/test-workspace?token="+validToken, nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !handlerCalled {
		t.Error("Handler should have been called with valid token in query param")
	}
	if rr.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rr.Code)
	}
}
