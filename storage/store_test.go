package storage

import (
	"labops/model"
	"path/filepath"
	"testing"
	"time"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "lab.db")
	s, e := Open(p)
	if e != nil {
		t.Fatal(e)
	}
	r := model.NewRecord("r1", "e1", "persist", time.Now())
	if e = s.SaveRecord(r); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.GetRecord("r1")
	if e != nil || got.Title != "persist" {
		t.Fatalf("%v %#v", e, got)
	}
}
