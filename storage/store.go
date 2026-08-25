package storage

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"path/filepath"
	"sync"
	"time"
)

var buckets = [][]byte{[]byte("records"), []byte("users"), []byte("events"), []byte("audits"), []byte("equipment"), []byte("assignments"), []byte("plans")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func encode(v any) ([]byte, error) { return json.Marshal(v) }
func decode(b []byte, v any) error { return json.Unmarshal(b, v) }
func (s *Store) put(bucket, key string, v any) error {
	b, err := encode(v)
	if err != nil {
		return err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), b) })
}
func (s *Store) get(bucket, key string, v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return fmt.Errorf("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if b == nil {
			return bbolt.ErrBucketNotFound
		}
		return decode(b, v)
	})
}
func (s *Store) delete(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete([]byte(key)) })
}
