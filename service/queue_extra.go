package service

import (
	"fmt"
	"time"
)

func (q *Queue) Promote(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	if !ok {
		return fmt.Errorf("work item not found")
	}
	if item.Status != "queued" {
		return fmt.Errorf("only queued work promotes")
	}
	item.Priority = "urgent"
	q.items[id] = item
	return nil
}
func (q *Queue) Delay(id string, until time.Time) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	if !ok {
		return fmt.Errorf("work item not found")
	}
	if until.Before(item.CreatedAt) {
		return fmt.Errorf("delay must be after creation")
	}
	item.Status = "delayed"
	item.FinishedAt = until
	q.items[id] = item
	return nil
}
func (q *Queue) ReleaseDelayed(now time.Time) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	count := 0
	for id, item := range q.items {
		if item.Status == "delayed" && !item.FinishedAt.After(now) {
			item.Status = "queued"
			q.items[id] = item
			count++
		}
	}
	return count
}
func (q *Queue) RemoveFinished(before time.Time) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	count := 0
	for id, item := range q.items {
		if item.Status == "done" && item.FinishedAt.Before(before) {
			delete(q.items, id)
			count++
		}
	}
	return count
}
func (q *Queue) FailedAttempts() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	total := 0
	for _, item := range q.items {
		if item.Status == "retry" {
			total += item.Attempts
		}
	}
	return total
}
func (q *Queue) HasQueued(recordID string) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, item := range q.items {
		if item.RecordID == recordID && (item.Status == "queued" || item.Status == "running" || item.Status == "retry") {
			return true
		}
	}
	return false
}
func (q *Queue) CancelRecord(recordID string) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	count := 0
	for id, item := range q.items {
		if item.RecordID == recordID && item.Status != "done" && item.Status != "cancelled" {
			item.Status = "cancelled"
			q.items[id] = item
			count++
		}
	}
	return count
}
func (q *Queue) Oldest(status string) (WorkItem, bool) {
	items := q.List(status)
	if len(items) == 0 {
		return WorkItem{}, false
	}
	old := items[0]
	for _, item := range items[1:] {
		if item.CreatedAt.Before(old.CreatedAt) {
			old = item
		}
	}
	return old, true
}
