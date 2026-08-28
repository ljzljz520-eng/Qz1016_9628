package workflow

import (
	"labops/model"
	"labops/service"
	"labops/storage"
	"path/filepath"
	"testing"
	"time"
)

func setup(t *testing.T) (*Intake, *Processing) {
	s, _ := storage.Open(filepath.Join(t.TempDir(), "x"))
	t.Cleanup(func() { s.Close() })
	rs := service.NewRecordService(s)
	return NewIntake(rs), NewProcessing(rs)
}
func TestWorkflowOne(t *testing.T) {
	i, _ := setup(t)
	r := model.NewRecord("one", "e", "one", time.Now())
	if i.Receive(r, "u") != nil || i.Validate("one", "u") != nil || i.Save("one", "u") != nil {
		t.Fatal("workflow")
	}
	got, e := i.Display("one")
	if e != nil || got.Status != "processing" {
		t.Fatalf("expected current status, got %s", got.Status)
	}
}
func TestWorkflowTwo(t *testing.T) {
	i, p := setup(t)
	i.Receive(model.NewRecord("two", "e", "two", time.Now()), "u")
	i.Validate("two", "u")
	if p.Resolve("two", "u") == nil {
		t.Fatal("resolve should reject processing")
	}
}
func TestWorkflowThree(t *testing.T) {
	i, p := setup(t)
	i.Receive(model.NewRecord("three", "e", "three", time.Now()), "u")
	i.Validate("three", "u")
	i.Save("three", "u")
	if p.Resolve("three", "u") != nil {
		t.Fatal("resolve")
	}
}
