package storage

import (
	"go.etcd.io/bbolt"
	"labops/model"
	"sort"
)

func (s *Store) SaveRecord(r model.Record) error { return s.put("records", r.ID, r) }
func (s *Store) GetRecord(id string) (*model.Record, error) {
	var r model.Record
	if err := s.get("records", id, &r); err != nil {
		return nil, err
	}
	return &r, nil
}
func (s *Store) ListRecords() ([]model.Record, error) {
	out := []model.Record{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("records")).ForEach(func(_, v []byte) error {
			var r model.Record
			if err := decode(v, &r); err != nil {
				return err
			}
			out = append(out, r)
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].UpdatedAt.Before(out[j].UpdatedAt) })
	return out, err
}
func (s *Store) DeleteRecord(id string) error { return s.delete("records", id) }
