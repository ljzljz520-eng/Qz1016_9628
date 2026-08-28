package api

import (
	"labops/service"
	"labops/storage"
	"labops/workflow"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestHealth(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x"))
	defer s.Close()
	h := NewServer(service.NewRecordService(s), workflow.NewQuery(s))
	w := httptest.NewRecorder()
	h.Health(w, httptest.NewRequest("GET", "/health", nil))
	if w.Code != 200 {
		t.Fatal(w.Code)
	}
}
