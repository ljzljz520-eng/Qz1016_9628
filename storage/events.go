package storage

import (
	"go.etcd.io/bbolt"
	"labops/model"
)

func (s *Store) SaveEvent(e model.Event) error { return s.put("events", e.ID, e) }
func (s *Store) ListEvents(recordID string) ([]model.Event, error) {
	out := []model.Event{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("events")).ForEach(func(_, v []byte) error {
			var e model.Event
			if err := decode(v, &e); err != nil {
				return err
			}
			if recordID == "" || e.RecordID == recordID {
				out = append(out, e)
			}
			return nil
		})
	})
	return out, err
}
func (s *Store) SaveAudit(a model.Audit) error { return s.put("audits", a.ID, a) }
func (s *Store) ListAudits(recordID string) ([]model.Audit, error) {
	out := []model.Audit{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("audits")).ForEach(func(_, v []byte) error {
			var a model.Audit
			if err := decode(v, &a); err != nil {
				return err
			}
			if recordID == "" || a.RecordID == recordID {
				out = append(out, a)
			}
			return nil
		})
	})
	return out, err
}
