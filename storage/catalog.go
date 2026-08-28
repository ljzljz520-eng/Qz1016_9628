package storage

import (
	"encoding/json"
	"fmt"
	"go.etcd.io/bbolt"
	"labops/model"
	"sort"
	"time"
)

var planBucket = []byte("plans")

func (s *Store) ensurePlanBucket() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.db.Update(func(tx *bbolt.Tx) error { _, e := tx.CreateBucketIfNotExists(planBucket); return e })
}
func (s *Store) SavePlan(p model.MaintenancePlan) error {
	if e := p.Validate(); e != nil {
		return e
	}
	if e := s.ensurePlanBucket(); e != nil {
		return e
	}
	return s.put("plans", p.ID, p)
}
func (s *Store) GetPlan(id string) (*model.MaintenancePlan, error) {
	if e := s.ensurePlanBucket(); e != nil {
		return nil, e
	}
	var p model.MaintenancePlan
	if e := s.get("plans", id, &p); e != nil {
		return nil, e
	}
	return &p, nil
}
func (s *Store) ListPlans(equipment string) ([]model.MaintenancePlan, error) {
	if e := s.ensurePlanBucket(); e != nil {
		return nil, e
	}
	out := []model.MaintenancePlan{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	e := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket(planBucket).ForEach(func(_, v []byte) error {
			var p model.MaintenancePlan
			if x := json.Unmarshal(v, &p); x != nil {
				return x
			}
			if equipment == "" || p.EquipmentID == equipment {
				out = append(out, p)
			}
			return nil
		})
	})
	sort.Slice(out, func(i, j int) bool { return out[i].NextDue.Before(out[j].NextDue) })
	return out, e
}
func (s *Store) CompletePlan(id string, now time.Time) error {
	p, e := s.GetPlan(id)
	if e != nil {
		return e
	}
	*p = p.Schedule(now)
	return s.SavePlan(*p)
}
func (s *Store) DuePlans(now time.Time) ([]model.MaintenancePlan, error) {
	rows, e := s.ListPlans("")
	if e != nil {
		return nil, e
	}
	out := []model.MaintenancePlan{}
	for _, p := range rows {
		if p.Due(now) {
			out = append(out, p)
		}
	}
	return out, nil
}
func (s *Store) SaveAssignment(a model.Assignment) error {
	if !a.Valid() {
		return fmt.Errorf("invalid assignment")
	}
	return s.put("assignments", a.RecordID+"/"+a.UserID, a)
}
func (s *Store) GetAssignments(record string) ([]model.Assignment, error) {
	out := []model.Assignment{}
	s.mu.RLock()
	defer s.mu.RUnlock()
	return out, s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte("assignments"))
		if b == nil {
			return nil
		}
		return b.ForEach(func(_, v []byte) error {
			var a model.Assignment
			if e := json.Unmarshal(v, &a); e != nil {
				return e
			}
			if record == "" || a.RecordID == record {
				out = append(out, a)
			}
			return nil
		})
	})
}
