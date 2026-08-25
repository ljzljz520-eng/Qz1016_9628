package report

import (
	"fmt"
	"labops/model"
	"labops/workflow"
)

type Builder struct{ Query *workflow.Query }

func NewBuilder(q *workflow.Query) *Builder { return &Builder{Query: q} }
func (b *Builder) StatusSummary(status string) (string, error) {
	rows, e := b.Query.Search("", status)
	if e != nil {
		return "", e
	}
	return fmt.Sprintf("%s:%d", status, len(rows)), nil
}
func FormatRecord(r model.Record) string { return fmt.Sprintf("%s [%s] %s", r.ID, r.Status, r.Title) }
func Render(rows []model.Record) []string {
	out := make([]string, 0, len(rows))
	for _, r := range rows {
		out = append(out, FormatRecord(r))
	}
	return out
}
