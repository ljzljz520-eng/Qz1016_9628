package workflow

import (
	"fmt"
	"labops/model"
	"labops/service"
)

type Intake struct{ Records *service.RecordService }

func NewIntake(r *service.RecordService) *Intake { return &Intake{Records: r} }
func (w *Intake) Receive(r model.Record, actor string) error {
	if r.Status == "" {
		r.Status = "new"
	}
	return w.Records.Register(r, actor)
}
func (w *Intake) Validate(id, actor string) error {
	return w.Records.Transition(id, "validated", actor)
}
func (w *Intake) Save(id, actor string) error {
	if err := w.Records.Transition(id, "processing", actor); err != nil {
		return fmt.Errorf("save workflow: %w", err)
	}
	return nil
}
func (w *Intake) Display(id string) (model.Record, error) { return w.Records.Snapshot(id) }
