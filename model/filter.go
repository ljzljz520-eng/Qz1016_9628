package model

import "sort"

type RecordFilter struct {
	Term string

	Statuses []string

	Priority string

	Assignee string

	Equipment string

	Limit int

	Offset int
}

func (f RecordFilter) Match(r Record) bool {
	if f.Term != "" && !containsFold(r.Title, f.Term) && !containsFold(r.Description, f.Term) {
		return false
	}
	if f.Priority != "" && r.Priority != f.Priority {
		return false
	}
	if f.Assignee != "" && r.Assignee != f.Assignee {
		return false
	}
	if f.Equipment != "" && r.EquipmentID != f.Equipment {
		return false
	}
	if len(f.Statuses) > 0 && !stringIn(f.Statuses, r.Status) {
		return false
	}
	return true
}

func containsFold(value, term string) bool {
	value = lower(value)
	term = lower(term)
	for i := 0; i+len(term) <= len(value); i++ {
		if value[i:i+len(term)] == term {
			return true
		}
	}
	return false
}

func lower(v string) string {
	result := make([]byte, len(v))
	for i := range v {
		c := v[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func stringIn(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func SortByUpdated(records []Record, descending bool) []Record {
	result := append([]Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if descending {
			return result[i].UpdatedAt.After(result[j].UpdatedAt)
		}
		return result[i].UpdatedAt.Before(result[j].UpdatedAt)
	})
	return result
}

func Paginate(records []Record, offset, limit int) []Record {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = len(records)
	}
	if offset >= len(records) {
		return []Record{}
	}
	end := offset + limit
	if end > len(records) {
		end = len(records)
	}
	return append([]Record(nil), records[offset:end]...)
}

func StatusCounts(records []Record) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}

func PriorityCounts(records []Record) map[string]int {
	counts := make(map[string]int)
	for _, record := range records {
		priority := NormalizePriority(record.Priority)
		counts[priority]++
	}
	return counts
}

func ValidateFilter(f RecordFilter) error {
	if f.Limit < 0 || f.Offset < 0 {
		return errInvalidPage
	}
	if f.Priority != "" && NormalizePriority(f.Priority) != f.Priority {
		return errInvalidPriority
	}
	for _, status := range f.Statuses {
		if !ValidStatus(status) {
			return errInvalidStatus
		}
	}
	return nil
}
