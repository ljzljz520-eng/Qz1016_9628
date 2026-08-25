package service

import (
	"fmt"
	"labops/model"
	"sort"
	"sync"
	"time"
)

type WorkItem struct {
	ID         string
	RecordID   string
	Queue      string
	Priority   string
	CreatedAt  time.Time
	StartedAt  time.Time
	FinishedAt time.Time
	Attempts   int
	Status     string
}
type Queue struct {
	mu    sync.Mutex
	items map[string]WorkItem
	clock func() time.Time
}

func NewQueue(clock func() time.Time) *Queue {
	if clock == nil {
		clock = time.Now
	}
	return &Queue{items: map[string]WorkItem{}, clock: clock}
}
func (q *Queue) Enqueue(record model.Record) WorkItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	item := WorkItem{ID: fmt.Sprintf("work-%s-%d", record.ID, len(q.items)+1), RecordID: record.ID, Queue: "maintenance", Priority: model.NormalizePriority(record.Priority), CreatedAt: q.clock(), Status: "queued"}
	q.items[item.ID] = item
	return item
}
func (q *Queue) Start(id string) (WorkItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	if !ok {
		return WorkItem{}, fmt.Errorf("work item not found")
	}
	if item.Status != "queued" && item.Status != "retry" {
		return WorkItem{}, fmt.Errorf("work item not startable")
	}
	item.Status = "running"
	item.StartedAt = q.clock()
	item.Attempts++
	q.items[id] = item
	return item, nil
}
func (q *Queue) Finish(id string, success bool) (WorkItem, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	if !ok {
		return WorkItem{}, fmt.Errorf("work item not found")
	}
	if item.Status != "running" {
		return WorkItem{}, fmt.Errorf("work item not running")
	}
	if success {
		item.Status = "done"
	} else {
		item.Status = "retry"
	}
	item.FinishedAt = q.clock()
	q.items[id] = item
	return item, nil
}
func (q *Queue) Cancel(id string) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	if !ok {
		return fmt.Errorf("work item not found")
	}
	if item.Status == "done" {
		return fmt.Errorf("completed work cannot cancel")
	}
	item.Status = "cancelled"
	q.items[id] = item
	return nil
}
func (q *Queue) Get(id string) (WorkItem, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	item, ok := q.items[id]
	return item, ok
}
func (q *Queue) List(status string) []WorkItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]WorkItem, 0)
	for _, item := range q.items {
		if status == "" || item.Status == status {
			out = append(out, item)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if model.ComparePriority(out[i].Priority, out[j].Priority) == 0 {
			return out[i].CreatedAt.Before(out[j].CreatedAt)
		}
		return model.ComparePriority(out[i].Priority, out[j].Priority) > 0
	})
	return out
}
func (q *Queue) Size(status string) int { return len(q.List(status)) }
func (q *Queue) Retryable(id string) bool {
	item, ok := q.Get(id)
	return ok && item.Status == "retry" && item.Attempts < 3
}
func (q *Queue) Drain() []WorkItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]WorkItem, 0, len(q.items))
	for _, item := range q.items {
		out = append(out, item)
	}
	q.items = map[string]WorkItem{}
	return out
}
func (q *Queue) Snapshot() map[string]WorkItem {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make(map[string]WorkItem, len(q.items))
	for k, v := range q.items {
		out[k] = v
	}
	return out
}
