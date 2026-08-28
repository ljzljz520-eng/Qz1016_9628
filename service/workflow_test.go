package service

import (
	"testing"
	"time"

	"labops/model"
	"labops/storage"
)

// TestWorkflowOne mirrors the documented workflow contract:
// 接收运维资料 -> 校验记录字段 -> 保存记录与事件 -> 展示最新进度.
// The Display step must show the NEW progress status ("processing"), not the
// stale status ("new") captured when the record was first received.
func TestWorkflowOne(t *testing.T) {
	dir := t.TempDir()
	s, err := storage.Open(dir + "/labops.db")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	rs := NewRecordService(s)
	rs.Clock = func() time.Time { return time.Unix(1_700_000_000, 0) }

	id := "REC-1"
	rec := model.NewRecord(id, "EQ-9", "示波器定期校准", time.Unix(1_700_000_000, 0))

	// 1. 接收运维资料
	if err := rs.Register(rec, "lab"); err != nil {
		t.Fatalf("register: %v", err)
	}

	// 2. 校验记录字段
	if err := rs.Transition(id, "validated", "lab"); err != nil {
		t.Fatalf("validate: %v", err)
	}

	// 3. 保存记录与事件
	if err := rs.Transition(id, "processing", "lab"); err != nil {
		t.Fatalf("save: %v", err)
	}

	// 4. 展示最新进度 — must reflect the latest status, not the cached one.
	got, err := rs.Snapshot(id)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if got.Status != "processing" {
		t.Fatalf("expected latest progress status %q, got stale status %q", "processing", got.Status)
	}
	if got.Version < 3 {
		t.Fatalf("expected version bumped by transitions, got %d", got.Version)
	}
}

// TestAssignRefreshesCache ensures Assign also keeps the snapshot in sync,
// so the displayed record reflects the newly assigned user.
func TestAssignRefreshesCache(t *testing.T) {
	dir := t.TempDir()
	s, err := storage.Open(dir + "/labops.db")
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	rs := NewRecordService(s)
	rs.Clock = func() time.Time { return time.Unix(1_700_000_000, 0) }

	id := "REC-2"
	rec := model.NewRecord(id, "EQ-1", "离心机维护", time.Unix(1_700_000_000, 0))
	if err := rs.Register(rec, "lab"); err != nil {
		t.Fatal(err)
	}

	if err := rs.Assign(id, "engineer-7", "lab"); err != nil {
		t.Fatalf("assign: %v", err)
	}

	got, err := rs.Snapshot(id)
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if got.Assignee != "engineer-7" {
		t.Fatalf("expected assignee %q, got %q", "engineer-7", got.Assignee)
	}
}
