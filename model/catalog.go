package model

import "time"

type MaintenancePlan struct {
	ID            string
	EquipmentID   string
	Name          string
	IntervalDays  int
	LastCompleted time.Time
	NextDue       time.Time
	Owner         string
	Enabled       bool
}

func (p MaintenancePlan) Due(now time.Time) bool {
	if !p.Enabled {
		return false
	}
	return !p.NextDue.After(now)
}
func (p MaintenancePlan) Schedule(now time.Time) MaintenancePlan {
	if p.IntervalDays <= 0 {
		p.IntervalDays = 30
	}
	p.LastCompleted = now
	p.NextDue = now.AddDate(0, 0, p.IntervalDays)
	return p
}
func (p MaintenancePlan) Validate() error {
	if p.ID == "" || p.EquipmentID == "" || p.Name == "" {
		return errInvalidPlan
	}
	if p.IntervalDays < 1 {
		return errInvalidPlan
	}
	return nil
}

type Assignment struct {
	RecordID   string
	UserID     string
	Role       string
	AssignedAt time.Time
	Active     bool
}

func (a Assignment) Valid() bool        { return a.RecordID != "" && a.UserID != "" && a.Role != "" }
func (a Assignment) Expire() Assignment { a.Active = false; return a }

type DailySummary struct {
	Day      time.Time
	Created  int
	Resolved int
	Blocked  int
	Closed   int
}

func (s DailySummary) Total() int { return s.Created + s.Resolved + s.Blocked + s.Closed }
func (s DailySummary) CompletionRate() float64 {
	if s.Created == 0 {
		return 0
	}
	return float64(s.Resolved+s.Closed) / float64(s.Created)
}
