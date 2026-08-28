package workflow

import (
	"sort"
	"time"

	"labops/model"
)

type Metrics struct {
	Total      int
	Open       int
	Closed     int
	ByStatus   map[string]int
	ByPriority map[string]int
	Oldest     time.Time
	Newest     time.Time
	AverageAge time.Duration
}

func BuildMetrics(records []model.Record, now time.Time) Metrics {
	metrics := Metrics{ByStatus: make(map[string]int), ByPriority: make(map[string]int)}
	if len(records) == 0 {
		return metrics
	}
	ages := make([]time.Duration, 0, len(records))
	for _, record := range records {
		metrics.Total++
		metrics.ByStatus[record.Status]++
		metrics.ByPriority[model.NormalizePriority(record.Priority)]++
		if record.IsOpen() {
			metrics.Open++
		} else {
			metrics.Closed++
		}
		if metrics.Oldest.IsZero() || record.CreatedAt.Before(metrics.Oldest) {
			metrics.Oldest = record.CreatedAt
		}
		if record.UpdatedAt.After(metrics.Newest) {
			metrics.Newest = record.UpdatedAt
		}
		if now.After(record.CreatedAt) {
			ages = append(ages, now.Sub(record.CreatedAt))
		}
	}
	for _, age := range ages {
		metrics.AverageAge += age
	}
	if len(ages) > 0 {
		metrics.AverageAge /= time.Duration(len(ages))
	}
	return metrics
}

func TopAssignees(records []model.Record) []string {
	counts := make(map[string]int)
	for _, record := range records {
		if record.Assignee != "" {
			counts[record.Assignee]++
		}
	}
	type pair struct {
		name  string
		count int
	}
	pairs := make([]pair, 0, len(counts))
	for name, count := range counts {
		pairs = append(pairs, pair{name, count})
	}
	sort.Slice(pairs, func(i, j int) bool {
		if pairs[i].count == pairs[j].count {
			return pairs[i].name < pairs[j].name
		}
		return pairs[i].count > pairs[j].count
	})
	result := make([]string, len(pairs))
	for i := range pairs {
		result[i] = pairs[i].name
	}
	return result
}

func NeedsAttention(record model.Record) bool {
	if record.Status == "blocked" {
		return true
	}
	if model.IsHighPriority(record.Priority) && record.IsOpen() {
		return true
	}
	return false
}

func AttentionList(records []model.Record) []model.Record {
	result := make([]model.Record, 0)
	for _, record := range records {
		if NeedsAttention(record) {
			result = append(result, record)
		}
	}
	return model.SortByUpdated(result, true)
}
