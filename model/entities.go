package model

import "time"

type Record struct {
	ID, EquipmentID, Title, Description, Status, Priority, Assignee string
	CreatedAt, UpdatedAt                                            time.Time
	Version                                                         int
}
type User struct {
	ID, Name, Role, Email string
	Active                bool
	CreatedAt             time.Time
}
type Event struct {
	ID, RecordID, Kind, Actor, Detail string
	At                                time.Time
}
type Audit struct {
	ID, RecordID, Action, Actor, Before, After string
	At                                         time.Time
}
type Equipment struct {
	ID, Name, Category, Location, Serial string
	Active                               bool
}

func (r Record) IsOpen() bool { return r.Status != "closed" && r.Status != "archived" }
func (r Record) CanTransition(to string) bool {
	switch r.Status {
	case "new":
		return to == "validated" || to == "cancelled"
	case "validated":
		return to == "processing" || to == "cancelled"
	case "processing":
		return to == "resolved" || to == "blocked"
	case "blocked":
		return to == "processing" || to == "cancelled"
	case "resolved":
		return to == "closed"
	case "closed":
		return to == "archived"
	}
	return false
}
func (r Record) Touch(now time.Time) Record { r.UpdatedAt = now; r.Version++; return r }
func NewRecord(id, equipment, title string, now time.Time) Record {
	return Record{ID: id, EquipmentID: equipment, Title: title, Status: "new", Priority: "normal", CreatedAt: now, UpdatedAt: now, Version: 1}
}
