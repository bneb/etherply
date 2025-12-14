package store

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
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
	key := []byte(fmt.Sprintf("project:%s", p.ID))
	val, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("failed to marshal project: %w", err)
	}

	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set(key, val)
	})
}

// ListProjects retrieves all projects.
// Warning: Full scan. In production, use a prefix iterator with pagination.
func (b *BadgerStore) ListProjects() ([]Project, error) {
	var projects []Project

	err := b.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte("project:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				var p Project
				if err := json.Unmarshal(val, &p); err != nil {
					return err
				}
				projects = append(projects, p)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	return projects, err
}
