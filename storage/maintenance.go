package storage

import "time"

func (s *Store) PurgeBefore(cutoff time.Time) (int, error) {
	rows, e := s.ListRecords()
	if e != nil {
		return 0, e
	}
	n := 0
	for _, r := range rows {
		if r.UpdatedAt.Before(cutoff) && r.Status == "archived" {
			if s.DeleteRecord(r.ID) == nil {
				n++
			}
		}
	}
	return n, nil
}
func (s *Store) Health() error { _, e := s.ListRecords(); return e }
