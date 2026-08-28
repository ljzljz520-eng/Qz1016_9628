package storage

import (
	"sort"
	"strings"

	"go.etcd.io/bbolt"
	"labops/model"
)

type RecordPage struct {
	Items  []model.Record
	Total  int
	Offset int
	Limit  int
}

func (s *Store) FindRecords(filter model.RecordFilter) (RecordPage, error) {
	if err := model.ValidateFilter(filter); err != nil {
		return RecordPage{}, err
	}
	rows, err := s.ListRecords()
	if err != nil {
		return RecordPage{}, err
	}
	matched := make([]model.Record, 0, len(rows))
	for _, row := range rows {
		if filter.Match(row) {
			matched = append(matched, row)
		}
	}
	sort.SliceStable(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	page := model.Paginate(matched, filter.Offset, filter.Limit)
	return RecordPage{Items: page, Total: len(matched), Offset: filter.Offset, Limit: filter.Limit}, nil
}

func (s *Store) FindByTitle(term string) ([]model.Record, error) {
	return s.FindByFilter(model.RecordFilter{Term: strings.TrimSpace(term)})
}

func (s *Store) FindByFilter(filter model.RecordFilter) ([]model.Record, error) {
	page, err := s.FindRecords(filter)
	if err != nil {
		return nil, err
	}
	return page.Items, nil
}

func (s *Store) CountByStatus() (map[string]int, error) {
	rows, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	return model.StatusCounts(rows), nil
}

func (s *Store) CountByPriority() (map[string]int, error) {
	rows, err := s.ListRecords()
	if err != nil {
		return nil, err
	}
	return model.PriorityCounts(rows), nil
}

func (s *Store) RecordExists(id string) (bool, error) {
	record, err := s.GetRecord(id)
	if err != nil {
		return false, nil
	}
	return record != nil, nil
}

func (s *Store) CopyRecord(source, target string) error {
	record, err := s.GetRecord(source)
	if err != nil {
		return err
	}
	record.ID = target
	return s.SaveRecord(*record)
}

func (s *Store) ReplaceStatus(id, status string) error {
	record, err := s.GetRecord(id)
	if err != nil {
		return err
	}
	if !model.ValidStatus(status) {
		return model.ValidateRecord(model.Record{ID: record.ID, EquipmentID: record.EquipmentID, Title: record.Title, Status: status})
	}
	record.Status = status
	return s.SaveRecord(*record)
}

func (s *Store) TouchRecord(id string) error {
	record, err := s.GetRecord(id)
	if err != nil {
		return err
	}
	record.UpdatedAt = record.UpdatedAt.Add(1)
	record.Version++
	return s.SaveRecord(*record)
}

func (s *Store) AllEquipment() ([]model.Equipment, error) {
	result := make([]model.Equipment, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte("equipment")).ForEach(func(_, value []byte) error {
			var item model.Equipment
			if err := decode(value, &item); err != nil {
				return err
			}
			result = append(result, item)
			return nil
		})
	})
	return result, err
}
