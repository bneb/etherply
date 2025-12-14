package sync

import (
	"encoding/json"
	"time"
)

// ServerAuthStrategy implements SyncStrategy where the server state is authoritative.
// In merge conflicts, local (server) state always wins.
// This is suitable for game state, admin-controlled settings, or non-collaborative data.
type ServerAuthStrategy struct{}

// NewServerAuthStrategy creates a new server-authoritative sync strategy.
func NewServerAuthStrategy() *ServerAuthStrategy {
	return &ServerAuthStrategy{}
}

// Name returns the strategy identifier.
func (s *ServerAuthStrategy) Name() string {
	return string(StrategyServerAuthoritative)
}

// ProcessWrite applies a mutation directly (server always succeeds).
func (s *ServerAuthStrategy) ProcessWrite(current []byte, key string, value interface{}, ts time.Time) ([]byte, error) {
	doc := s.loadOrCreate(current)
	doc[key] = value
	return json.Marshal(doc)
}

// Merge keeps local (server) state, ignoring remote.
func (s *ServerAuthStrategy) Merge(local, remote []byte) ([]byte, error) {
	if len(local) == 0 {
		return remote, nil
	}
	// Server wins - return local unchanged
	return local, nil
}

// GetState returns the document as-is.
func (s *ServerAuthStrategy) GetState(doc []byte) (map[string]interface{}, error) {
	return s.loadOrCreate(doc), nil
}

// GetHeads returns empty - no version tracking.
func (s *ServerAuthStrategy) GetHeads(doc []byte) ([]string, error) {
	return []string{}, nil
}

// GetChanges returns full document (no incremental sync).
func (s *ServerAuthStrategy) GetChanges(doc []byte, since []string) ([]byte, error) {
	return doc, nil
}

// GetHistory returns empty - no history tracking.
func (s *ServerAuthStrategy) GetHistory(doc []byte) ([]Change, error) {
	return []Change{}, nil
}

func (s *ServerAuthStrategy) loadOrCreate(data []byte) map[string]interface{} {
	if len(data) == 0 {
		return make(map[string]interface{})
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return make(map[string]interface{})
	}
	return doc
}
