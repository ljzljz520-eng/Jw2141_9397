package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"go.etcd.io/bbolt"
)

var bucketNames = [][]byte{
	[]byte("farmers"),
	[]byte("visits"),
	[]byte("reviews"),
	[]byte("notes"),
	[]byte("archives"),
}

type Database struct {
	db *bbolt.DB
}

func Open(path string) (*Database, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bbolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	store := &Database{db: db}
	if err := store.createBuckets(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Database) createBuckets() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return fmt.Errorf("create bucket %s: %w", name, err)
			}
		}
		return nil
	})
}

func (s *Database) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Database) Put(bucket string, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode %s: %w", key, err)
	}
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucketRef := tx.Bucket([]byte(bucket))
		if bucketRef == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return bucketRef.Put([]byte(key), data)
	})
}

func (s *Database) Get(bucket string, key string, target any) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		bucketRef := tx.Bucket([]byte(bucket))
		if bucketRef == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		data := bucketRef.Get([]byte(key))
		if data == nil {
			return ErrNotFound
		}
		return json.Unmarshal(append([]byte(nil), data...), target)
	})
}

func (s *Database) Delete(bucket string, key string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucketRef := tx.Bucket([]byte(bucket))
		if bucketRef == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return bucketRef.Delete([]byte(key))
	})
}

func (s *Database) Keys(bucket string) ([]string, error) {
	keys := make([]string, 0)
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucketRef := tx.Bucket([]byte(bucket))
		if bucketRef == nil {
			return fmt.Errorf("unknown bucket %s", bucket)
		}
		return bucketRef.ForEach(func(key, _ []byte) error {
			keys = append(keys, string(key))
			return nil
		})
	})
	sort.Strings(keys)
	return keys, err
}

func (s *Database) Count(bucket string) (int, error) {
	keys, err := s.Keys(bucket)
	return len(keys), err
}

var ErrNotFound = errors.New("record not found")
