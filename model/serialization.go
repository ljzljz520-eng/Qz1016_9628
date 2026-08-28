package model

import (
	"encoding/json"
	"fmt"
	"time"
)

func EncodeRecord(r Record) ([]byte, error) { return json.Marshal(r) }
func DecodeRecord(data []byte) (Record, error) {
	var r Record
	if e := json.Unmarshal(data, &r); e != nil {
		return Record{}, e
	}
	return r, nil
}
func EncodeUser(u User) ([]byte, error) { return json.Marshal(u) }
func DecodeUser(data []byte) (User, error) {
	var u User
	if e := json.Unmarshal(data, &u); e != nil {
		return User{}, e
	}
	return u, nil
}
func EncodeEvent(e Event) ([]byte, error) { return json.Marshal(e) }
func DecodeEvent(data []byte) (Event, error) {
	var e Event
	if x := json.Unmarshal(data, &e); x != nil {
		return Event{}, x
	}
	return e, nil
}
func EncodeAudit(a Audit) ([]byte, error) { return json.Marshal(a) }
func DecodeAudit(data []byte) (Audit, error) {
	var a Audit
	if e := json.Unmarshal(data, &a); e != nil {
		return Audit{}, e
	}
	return a, nil
}
func RecordMap(r Record) map[string]string {
	return map[string]string{"id": r.ID, "equipment_id": r.EquipmentID, "title": r.Title, "description": r.Description, "status": r.Status, "priority": r.Priority, "assignee": r.Assignee, "created_at": r.CreatedAt.Format(time.RFC3339Nano), "updated_at": r.UpdatedAt.Format(time.RFC3339Nano), "version": fmt.Sprint(r.Version)}
}
func RecordFromMap(values map[string]string) (Record, error) {
	created, e := time.Parse(time.RFC3339Nano, values["created_at"])
	if e != nil {
		return Record{}, e
	}
	updated, e := time.Parse(time.RFC3339Nano, values["updated_at"])
	if e != nil {
		return Record{}, e
	}
	var version int
	if _, e = fmt.Sscan(values["version"], &version); e != nil {
		return Record{}, e
	}
	return Record{ID: values["id"], EquipmentID: values["equipment_id"], Title: values["title"], Description: values["description"], Status: values["status"], Priority: values["priority"], Assignee: values["assignee"], CreatedAt: created, UpdatedAt: updated, Version: version}, nil
}
func MergeRecord(base, patch Record) Record {
	if patch.Title != "" {
		base.Title = patch.Title
	}
	if patch.Description != "" {
		base.Description = patch.Description
	}
	if patch.Priority != "" {
		base.Priority = NormalizePriority(patch.Priority)
	}
	if patch.Assignee != "" {
		base.Assignee = patch.Assignee
	}
	if patch.Status != "" {
		base.Status = patch.Status
	}
	return base
}
func NormalizeRecord(r Record) Record {
	r.ID = CanonicalTitle(r.ID)
	r.Title = CanonicalTitle(r.Title)
	r.Priority = NormalizePriority(r.Priority)
	if r.Status == "" {
		r.Status = "new"
	}
	return r
}
func EqualRecord(a, b Record) bool {
	return a.ID == b.ID && a.EquipmentID == b.EquipmentID && a.Title == b.Title && a.Description == b.Description && a.Status == b.Status && a.Priority == b.Priority && a.Assignee == b.Assignee && a.Version == b.Version
}
func RecordAge(r Record, now time.Time) time.Duration {
	if now.Before(r.CreatedAt) {
		return 0
	}
	return now.Sub(r.CreatedAt)
}
func IsRecent(r Record, now time.Time, window time.Duration) bool { return RecordAge(r, now) <= window }
func EventLabel(e Event) string {
	if e.Kind == "created" {
		return "记录创建"
	}
	if e.Kind == "transition" {
		return "状态变更"
	}
	return e.Kind
}
