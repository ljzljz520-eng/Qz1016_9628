package audit

import (
	"labops/model"
	"labops/storage"
)

type Reporter struct{ Store *storage.Store }

func NewReporter(s *storage.Store) *Reporter                 { return &Reporter{Store: s} }
func (r *Reporter) History(id string) ([]model.Audit, error) { return r.Store.ListAudits(id) }
func (r *Reporter) Count(id string) (int, error)             { a, e := r.History(id); return len(a), e }
