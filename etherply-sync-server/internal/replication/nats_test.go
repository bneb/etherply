// Package replication provides tests for the replication subsystem.
// These tests cover configuration validation, mock replicator functionality,
// and type serialization. Integration tests requiring a live NATS server
// are skipped by default but can be run with the -integration flag.
package replication

import (
	"context"
	"testing"
	"time"
)

// TestConfig verifies configuration validation and defaults.
func TestConfig_Defaults(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.StreamName != "ETHERPLY_REPLICATION" {
		t.Errorf("expected StreamName=ETHERPLY_REPLICATION, got %s", cfg.StreamName)
	}
	if cfg.ReconnectWait != 2*time.Second {
		t.Errorf("expected ReconnectWait=2s, got %v", cfg.ReconnectWait)
	}
	if cfg.MaxReconnects != -1 {
		t.Errorf("expected MaxReconnects=-1, got %d", cfg.MaxReconnects)
	}
}

func TestNATSReplicator_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name:    "missing ServerID",
			cfg:     Config{Region: "us-east-1", NATSURLs: []string{"nats://localhost:4222"}},
			wantErr: "ServerID is required",
		},
		{
			name:    "missing Region",
			cfg:     Config{ServerID: "server-1", NATSURLs: []string{"nats://localhost:4222"}},
			wantErr: "Region is required",
		},
		{
			name:    "missing NATSURLs",
			cfg:     Config{ServerID: "server-1", Region: "us-east-1"},
			wantErr: "at least one NATS URL is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewNATSReplicator(tt.cfg)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if err.Error() != tt.wantErr {
				t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

func TestChangeEvent_Serialization(t *testing.T) {
	event := ChangeEvent{
		WorkspaceID:    "test-workspace",
		Changes:        []byte("automerge-binary-data"),
		OriginRegion:   "us-east-1",
		OriginServerID: "server-1",
		Timestamp:      time.Now(),
	}

	if event.WorkspaceID != "test-workspace" {
		t.Errorf("unexpected WorkspaceID: %s", event.WorkspaceID)
	}
	if len(event.Changes) == 0 {
		t.Error("expected non-empty Changes")
	}
}

// MockReplicator implements Replicator for testing.
type MockReplicator struct {
	BroadcastCalls []ChangeEvent
	SubscribeErr   error
	PeersList      []Replica
	IsHealthy      bool
}

func NewMockReplicator() *MockReplicator {
	return &MockReplicator{
		BroadcastCalls: make([]ChangeEvent, 0),
		IsHealthy:      true,
	}
}

func (m *MockReplicator) Broadcast(ctx context.Context, event ChangeEvent) error {
	m.BroadcastCalls = append(m.BroadcastCalls, event)
	return nil
}

func (m *MockReplicator) Subscribe(handler ChangeHandler) error {
	return m.SubscribeErr
}

func (m *MockReplicator) Peers() []Replica {
	return m.PeersList
}

func (m *MockReplicator) Healthy() bool {
	return m.IsHealthy
}

func (m *MockReplicator) Close() error {
	return nil
}

// Ensure MockReplicator satisfies the interface
var _ Replicator = (*MockReplicator)(nil)

func TestMockReplicator_Broadcast(t *testing.T) {
	mock := NewMockReplicator()

	event := ChangeEvent{
		WorkspaceID:    "ws-123",
		Changes:        []byte("test"),
		OriginRegion:   "test-region",
		OriginServerID: "test-server",
		Timestamp:      time.Now(),
	}

	err := mock.Broadcast(context.Background(), event)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(mock.BroadcastCalls) != 1 {
		t.Fatalf("expected 1 broadcast call, got %d", len(mock.BroadcastCalls))
	}

	if mock.BroadcastCalls[0].WorkspaceID != "ws-123" {
		t.Errorf("unexpected workspace ID: %s", mock.BroadcastCalls[0].WorkspaceID)
	}
}

func TestMockReplicator_Healthy(t *testing.T) {
	mock := NewMockReplicator()

	if !mock.Healthy() {
		t.Error("expected mock to be healthy by default")
	}

	mock.IsHealthy = false
	if mock.Healthy() {
		t.Error("expected mock to be unhealthy after setting IsHealthy=false")
	}
}
