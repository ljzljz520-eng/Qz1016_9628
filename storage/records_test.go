package storage

import (
	"labops/model"
	"path/filepath"
	"testing"
	"time"
)

func TestRecordStorageList(t *testing.T) {
	s, e := Open(filepath.Join(t.TempDir(), "x"))
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	s.SaveRecord(model.NewRecord("a", "e", "a", time.Now()))
	s.SaveRecord(model.NewRecord("b", "e", "b", time.Now()))
	rows, e := s.ListRecords()
	if e != nil || len(rows) != 2 {
		t.Fatal(e, len(rows))
	}
}
