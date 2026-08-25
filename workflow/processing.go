package workflow

import (
	"fmt"
	"labops/service"
)

type Processing struct{ Records *service.RecordService }

func NewProcessing(r *service.RecordService) *Processing { return &Processing{Records: r} }
func (w *Processing) Start(id, actor string) error {
	return w.Records.Transition(id, "processing", actor)
}
func (w *Processing) Resolve(id, actor string) error {
	return w.Records.Transition(id, "resolved", actor)
}
func (w *Processing) Block(id, actor string) error { return w.Records.Transition(id, "blocked", actor) }
func (w *Processing) Archive(id, actor string) error {
	if err := w.Records.Transition(id, "closed", actor); err != nil {
		return err
	}
	if err := w.Records.Transition(id, "archived", actor); err != nil {
		return fmt.Errorf("archive: %w", err)
	}
	return nil
}
