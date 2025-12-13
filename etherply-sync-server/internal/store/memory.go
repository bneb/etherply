package store

import (
	"sync"
)

// StateStore interface captures the persistence requirements.
// For MVP, this is implemented as an in-memory map.
// In the future, this will be swapped for FoundationDB/CockroachDB.
type StateStore interface {
	Get(workspaceID string, key string) (interface{}, bool)
	Set(workspaceID string, key string, value interface{}) error
	GetAll(workspaceID string) (map[string]interface{}, error)
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]map[string]interface{} // workspaceID -> key -> value
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]map[string]interface{}),
	}
}

func (s *MemoryStore) Get(workspaceID string, key string) (interface{}, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[workspaceID]
	if !ok {
		return nil, false
	}
	val, ok := workspace[key]
	return val, ok
}

func (s *MemoryStore) Set(workspaceID string, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	workspace, ok := s.data[workspaceID]
	if !ok {
		workspace = make(map[string]interface{})
		s.data[workspaceID] = workspace
	}
	workspace[key] = value
	return nil
}

func (s *MemoryStore) GetAll(workspaceID string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[workspaceID]
	if !ok {
		// Return empty map if no data exists yet
		return make(map[string]interface{}), nil
	}

	// Return a copy to be safe
	result := make(map[string]interface{})
	for k, v := range workspace {
		result[k] = v
	}
	return result, nil
}
