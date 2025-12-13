package store

import (
	"bytes"
	"encoding/gob"
	"fmt"
	badger "github.com/dgraph-io/badger/v4"
)

// BadgerStore implements Store using BadgerDB v4.
type BadgerStore struct {
	db *badger.DB
}

// NewBadgerStore creates a new BadgerStore backed by the directory at path.
func NewBadgerStore(path string) (*BadgerStore, error) {
	opts := badger.DefaultOptions(path)
	// Reduce logging noise
	opts.Logger = nil 
	
	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger db: %w", err)
	}

	return &BadgerStore{
		db: db,
	}, nil
}

func (s *BadgerStore) Get(workspaceID, key string) (interface{}, bool, error) {
	var val []byte
	dbKey := makeKey(workspaceID, key)

	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(dbKey)
		if err == badger.ErrKeyNotFound {
			return err
		}
		if err != nil {
			return err
		}
		val, err = item.ValueCopy(nil)
		return err
	})

	if err == badger.ErrKeyNotFound {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	decoded, err := decode(val)
	if err != nil {
		return nil, true, fmt.Errorf("failed to decode value: %w", err)
	}

	return decoded, true, nil
}

func (s *BadgerStore) Set(workspaceID, key string, value interface{}) error {
	dbKey := makeKey(workspaceID, key)
	valBytes, err := encode(value)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(dbKey, valBytes)
	})
}

func (s *BadgerStore) GetAll(workspaceID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	prefix := []byte(workspaceID + ":")

	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = true
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			
			// Extract original key from composite key "workspace:key"
			// k is []byte
			keyStr := string(k)
			// primitive split, ideally robust
			originalKey := keyStr[len(prefix):]

			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			decoded, err := decode(valCopy)
			if err != nil {
				return fmt.Errorf("failed to decode value for key %s: %w", keyStr, err)
			}
			result[originalKey] = decoded
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *BadgerStore) Close() error {
	return s.db.Close()
}

// Helpers

func makeKey(workspaceID, key string) []byte {
	return []byte(fmt.Sprintf("%s:%s", workspaceID, key))
}

func encode(val interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(&val); err != nil { // Encode pointer to interface to capture concrete type
		return nil, err
	}
	return buf.Bytes(), nil
}

func decode(data []byte) (interface{}, error) {
	var buf bytes.Buffer
	buf.Write(data)
	dec := gob.NewDecoder(&buf)
	var val interface{}
	if err := dec.Decode(&val); err != nil {
		return nil, err
	}
	return val, nil
}

// Ensure interface satisfaction
var _ Store = (*BadgerStore)(nil)
