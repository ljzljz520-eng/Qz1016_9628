package workflow

import (
	"labops/model"
	"labops/storage"
	"strings"
)

type Query struct{ Store *storage.Store }

func NewQuery(s *storage.Store) *Query { return &Query{Store: s} }
func (q *Query) Search(term, status string) ([]model.Record, error) {
	rows, e := q.Store.ListRecords()
	if e != nil {
		return nil, e
	}
	out := make([]model.Record, 0)
	for _, r := range rows {
		if (term == "" || strings.Contains(strings.ToLower(r.Title), strings.ToLower(term))) && (status == "" || r.Status == status) {
			out = append(out, r)
		}
	}
	return out, nil
}
func (q *Query) Timeline(id string) ([]model.Event, error) { return q.Store.ListEvents(id) }
