package store

import (
	"sync"
)

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]map[string]interface{} // workspaceID -> key -> value
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]map[string]interface{}),
	}
}

func (s *MemoryStore) Get(namespace string, key string) (interface{}, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[namespace]
	if !ok {
		return nil, false, nil
	}
	val, ok := workspace[key]
	return val, ok, nil
}

func (s *MemoryStore) Set(namespace string, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	workspace, ok := s.data[namespace]
	if !ok {
		workspace = make(map[string]interface{})
		s.data[namespace] = workspace
	}
	workspace[key] = value
	return nil
}

func (s *MemoryStore) GetAll(namespace string) (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workspace, ok := s.data[namespace]
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

// Close is a no-op for MemoryStore but provided for interface consistency
// with DiskStore. This allows both stores to be used interchangeably.
func (s *MemoryStore) Close() error {
	return nil
}

func (s *MemoryStore) Stats() (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	totalKeys := 0
	for _, ws := range s.data {
		totalKeys += len(ws)
	}

	return map[string]interface{}{
		"workspaces": len(s.data),
		"keys":       totalKeys,
	}, nil
}

// Ping checks if the store is healthy. For MemoryStore, this always succeeds.
func (s *MemoryStore) Ping() error {
	return nil
}

// Ensure interface satisfaction
var _ Store = (*MemoryStore)(nil)
