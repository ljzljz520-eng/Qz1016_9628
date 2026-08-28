package service

import (
	"fmt"
	"labops/model"
	"labops/storage"
	"sync"
	"time"
)

type RecordService struct {
	Store *storage.Store
	Clock func() time.Time
	mu    sync.RWMutex
	cache map[string]model.Record
}

func NewRecordService(s *storage.Store) *RecordService {
	return &RecordService{Store: s, Clock: time.Now, cache: map[string]model.Record{}}
}
func (s *RecordService) Register(r model.Record, actor string) error {
	if err := model.ValidateRecord(r); err != nil {
		return err
	}
	if r.Status == "" {
		r.Status = "new"
	}
	if err := s.Store.SaveRecord(r); err != nil {
		return err
	}
	s.mu.Lock()
	s.cache[r.ID] = r
	s.mu.Unlock()
	return s.Store.SaveEvent(model.Event{ID: r.ID + "-created", RecordID: r.ID, Kind: "created", Actor: actor, At: s.Clock()})
}
func (s *RecordService) Transition(id, to, actor string) error {
	r, err := s.Store.GetRecord(id)
	if err != nil {
		return err
	}
	if !r.CanTransition(to) {
		return fmt.Errorf("transition %s to %s denied", r.Status, to)
	}
	before := r.Status
	r.Status = to
	*r = r.Touch(s.Clock())
	if err := s.Store.SaveRecord(*r); err != nil {
		return err
	}
	_ = s.Store.SaveAudit(model.Audit{ID: fmt.Sprintf("%s-%d", id, r.Version), RecordID: id, Action: "transition", Actor: actor, Before: before, After: to, At: s.Clock()})
	return nil
}
func (s *RecordService) Assign(id, user, actor string) error {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return e
	}
	r.Assignee = user
	*r = r.Touch(s.Clock())
	return s.Store.SaveRecord(*r)
}
func (s *RecordService) Snapshot(id string) (model.Record, error) {
	s.mu.RLock()
	cached, ok := s.cache[id]
	s.mu.RUnlock()
	if ok {
		return cached, nil
	}
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return model.Record{}, e
	}
	s.mu.Lock()
	s.cache[id] = *r
	s.mu.Unlock()
	return *r, nil
}
func (s *RecordService) LookupOrNil(id string) (*model.Record, error) {
	r, e := s.Store.GetRecord(id)
	if e != nil {
		return nil, nil
	}
	return r, nil
}
