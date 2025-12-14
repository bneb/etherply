package store

// Store defines the interface for persistence layers.
// It decouples the CRDT engine from the underlying storage technology (Badger, Postgres, etc).
type Store interface {
	// Get retrieves a value from the store.
	// Returns value, found boolean, and error.
	Get(namespace, key string) (interface{}, bool, error)

	// Set writes a value to the store.
	Set(namespace, key string, value interface{}) error

	// GetAll retrieves all key-value pairs for a given namespace.
	GetAll(namespace string) (map[string]interface{}, error)

	// Close cleans up resources.
	Close() error

	// Stats returns storage metrics (e.g., number of keys).
	Stats() (map[string]interface{}, error)

	// Ping checks if the store is healthy and accessible.
	// Returns nil if healthy, error otherwise.
	Ping() error
}
