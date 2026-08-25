package service

import (
	"labops/model"
	"labops/storage"
	"path/filepath"
	"testing"
	"time"
)

func TestRecordServiceTransition(t *testing.T) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x"))
	defer s.Close()
	rs := NewRecordService(s)
	r := model.NewRecord("r", "e", "x", time.Now())
	if rs.Register(r, "u") != nil {
		t.Fatal("register")
	}
	if rs.Transition("r", "validated", "u") != nil {
		t.Fatal("transition")
	}
}
