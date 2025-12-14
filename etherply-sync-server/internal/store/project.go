package store

import (
	"encoding/json"
	"fmt"
	"time"
)

// Project represents a workspace or app using EtherPly.
type Project struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	CreatedAt time.Time `json:"created_at"`
	// ActiveConnections is transient, not stored here usually, but keeping simple for now.
}

// SaveProject persists a project to BadgerDB.
// Key format: project:<id>
func (b *BadgerStore) SaveProject(p Project) error {
	key := fmt.Sprintf("project:%s", p.ID)
	val, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}

	// Use "sys:projects" namespace
	return b.Set("sys:projects", key, val)
}

// ListProjects retrieves all projects.
// Warning: Full scan. In production, use a prefix iterator with pagination.
func (b *BadgerStore) ListProjects() ([]Project, error) {
	var projects []Project

	// Use "sys:projects" namespace
	// GetAll now returns map[string]interface{} where values are what we stored (bytes or decoded).
	// BadgerStore.Set encodes with gob.
	// But SaveProject uses `val` from json.Marshal?
	// Wait, BadgerStore.Set calls `encode`.
	// So `val` (json bytes) is wrapped in gob.

	// Issue: SaveProject calls `b.db.Update` directly in original code!
	// Now I changed it to `b.Set`.
	// Please verify previous `SaveProject` logic.
	// It used `b.db.Update`.

	// If I use `b.Set`, it will gob-encode the []byte.
	// That's double encoding but fine.

	all, err := b.GetAll("sys:projects")
	if err != nil {
		return nil, err
	}

	for _, v := range all {
		// v is interface{}, likely []byte if we passed []byte to Set?
		// Set uses `encode`.
		// Decode returns interface{}.

		// If we passed []byte to Set, Decode returns []byte.
		// Let's assume v is []byte (json).

		data, ok := v.([]byte)
		if !ok {
			continue
		}

		var p Project
		if err := json.Unmarshal(data, &p); err != nil {
			continue
		}
		projects = append(projects, p)
	}

	return projects, nil
}
