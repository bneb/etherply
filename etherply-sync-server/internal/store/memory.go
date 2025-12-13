package store

import (
	"sync"
)

// StateStore defines the persistence contract for workspace state storage.
// Implementations must be thread-safe as they may be accessed concurrently.
//
// For MVP, this is implemented with an in-memory map (MemoryStore) or
// an Append-Only File (DiskStore). In production, this interface would
// be implemented by FoundationDB or CockroachDB adapters.
type StateStore interface {
	// Get retrieves a single value by workspace and key.
	// Returns (nil, false) if the key doesn't exist.
	Get(workspaceID string, key string) (interface{}, bool)

	// Set stores a value for the given workspace and key.
	// For persistent stores, this should also write to durable storage.
	Set(workspaceID string, key string, value interface{}) error

	// GetAll returns all key-value pairs for a workspace.
	// Returns an empty map (not nil) if the workspace has no data.
	GetAll(workspaceID string) (map[string]interface{}, error)

	// Close releases any resources held by the store (file handles, connections).
	// After Close is called, the store should not be used.
	Close() error
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

// Close is a no-op for MemoryStore but provided for interface consistency
// with DiskStore. This allows both stores to be used interchangeably.
func (s *MemoryStore) Close() error {
	return nil
}
