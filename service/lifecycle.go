package service

import (
	"fmt"
	"labops/audit"
	"labops/model"
	"labops/storage"
	"time"
)

type Lifecycle struct {
	Store    *storage.Store
	Records  *RecordService
	Notifier *Notifier
	Clock    func() time.Time
}

func NewLifecycle(s *storage.Store, r *RecordService, n *Notifier) *Lifecycle {
	return &Lifecycle{Store: s, Records: r, Notifier: n, Clock: time.Now}
}
func (l *Lifecycle) Submit(record model.Record, actor string) error {
	record.Status = "new"
	if e := l.Records.Register(record, actor); e != nil {
		return e
	}
	if l.Notifier != nil {
		l.Notifier.Notify(record, actor, "记录已接收")
	}
	return nil
}
func (l *Lifecycle) Review(id, role, actor, note string) error {
	r, e := l.Records.Snapshot(id)
	if e != nil {
		return e
	}
	p := audit.DefaultPolicy()
	if e = p.Check(role, r, "validated", note); e != nil {
		return e
	}
	if e = l.Records.Transition(id, "validated", actor); e != nil {
		return e
	}
	if l.Notifier != nil {
		r.Status = "validated"
		l.Notifier.NotifyTransition(r, actor)
	}
	return nil
}
func (l *Lifecycle) Begin(id, actor string) error {
	if e := l.Records.Transition(id, "processing", actor); e != nil {
		return e
	}
	r, e := l.Records.Snapshot(id)
	if e == nil && l.Notifier != nil {
		l.Notifier.Notify(*&r, actor, "开始处理")
	}
	return e
}
func (l *Lifecycle) Finish(id, actor string) error {
	if e := l.Records.Transition(id, "resolved", actor); e != nil {
		return e
	}
	r, e := l.Records.Snapshot(id)
	if e == nil && l.Notifier != nil {
		l.Notifier.Notify(r, actor, "处理完成")
	}
	return e
}
func (l *Lifecycle) Reject(id, actor string) error {
	r, e := l.Records.Snapshot(id)
	if e != nil {
		return e
	}
	if r.Status != "new" && r.Status != "validated" {
		return fmt.Errorf("cannot reject %s", r.Status)
	}
	return l.Records.Transition(id, "cancelled", actor)
}
func (l *Lifecycle) Close(id, actor string) error { return l.Records.Transition(id, "closed", actor) }
func (l *Lifecycle) Archive(id, actor string) error {
	return l.Records.Transition(id, "archived", actor)
}
func (l *Lifecycle) Reopen(id, actor string) error {
	r, e := l.Records.Snapshot(id)
	if e != nil {
		return e
	}
	if r.Status != "blocked" {
		return fmt.Errorf("only blocked records reopen")
	}
	return l.Records.Transition(id, "processing", actor)
}
func (l *Lifecycle) Assign(id, user, role, actor string) error {
	r, e := l.Records.Snapshot(id)
	if e != nil {
		return e
	}
	if !model.CanAssign(role, r.Status) {
		return fmt.Errorf("role cannot assign")
	}
	a := model.Assignment{RecordID: id, UserID: user, Role: role, AssignedAt: l.Clock(), Active: true}
	if e = l.Store.SaveAssignment(a); e != nil {
		return e
	}
	return l.Records.Assign(id, user, actor)
}
func (l *Lifecycle) Schedule(plan model.MaintenancePlan) error {
	if plan.NextDue.IsZero() {
		plan = plan.Schedule(l.Clock())
	}
	return l.Store.SavePlan(plan)
}
func (l *Lifecycle) RunDuePlans() ([]model.MaintenancePlan, error) {
	return l.Store.DuePlans(l.Clock())
}
