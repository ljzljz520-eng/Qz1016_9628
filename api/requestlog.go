package api

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RequestEntry struct {
	Method   string
	Path     string
	Status   int
	Started  time.Time
	Duration time.Duration
}
type RequestLog struct {
	mu      sync.Mutex
	entries []RequestEntry
	limit   int
}

func NewRequestLog(limit int) *RequestLog {
	if limit < 1 {
		limit = 100
	}
	return &RequestLog{entries: make([]RequestEntry, 0, limit), limit: limit}
}
func (l *RequestLog) Add(entry RequestEntry) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, entry)
	if len(l.entries) > l.limit {
		l.entries = l.entries[len(l.entries)-l.limit:]
	}
}
func (l *RequestLog) Entries() []RequestEntry {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]RequestEntry(nil), l.entries...)
}
func (l *RequestLog) Count() int { l.mu.Lock(); defer l.mu.Unlock(); return len(l.entries) }
func (l *RequestLog) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		started := time.Now()
		next.ServeHTTP(w, r)
		l.Add(RequestEntry{Method: r.Method, Path: r.URL.Path, Status: http.StatusOK, Started: started, Duration: time.Since(started)})
	})
}
func (l *RequestLog) Clear() { l.mu.Lock(); defer l.mu.Unlock(); l.entries = []RequestEntry{} }
func (l *RequestLog) Last() (RequestEntry, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.entries) == 0 {
		return RequestEntry{}, false
	}
	return l.entries[len(l.entries)-1], true
}
func (e RequestEntry) String() string {
	return fmt.Sprintf("%s %s %d %s", e.Method, e.Path, e.Status, e.Duration)
}
func (e RequestEntry) Slow(threshold time.Duration) bool { return e.Duration > threshold }

func (l *RequestLog) SlowEntries(threshold time.Duration) []RequestEntry {
	entries := l.Entries()
	result := make([]RequestEntry, 0)
	for _, entry := range entries {
		if entry.Slow(threshold) {
			result = append(result, entry)
		}
	}
	return result
}

func (l *RequestLog) ByPath(path string) []RequestEntry {
	entries := l.Entries()
	result := make([]RequestEntry, 0)
	for _, entry := range entries {
		if entry.Path == path {
			result = append(result, entry)
		}
	}
	return result
}

func (l *RequestLog) ErrorCount() int {
	entries := l.Entries()
	count := 0
	for _, entry := range entries {
		if entry.Status >= http.StatusBadRequest {
			count++
		}
	}
	return count
}

func (l *RequestLog) AverageDuration() time.Duration {
	entries := l.Entries()
	if len(entries) == 0 {
		return 0
	}
	var total time.Duration
	for _, entry := range entries {
		total += entry.Duration
	}
	return total / time.Duration(len(entries))
}
