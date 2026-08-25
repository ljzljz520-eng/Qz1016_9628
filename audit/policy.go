package audit

import (
	"fmt"
	"strings"

	"labops/model"
)

type Policy struct {
	AllowedRoles map[string][]string
	RequireNote  bool
}

func DefaultPolicy() Policy {
	return Policy{AllowedRoles: map[string][]string{
		"operator": {"new", "validated", "processing", "blocked"},
		"reviewer": {"validated", "resolved", "closed"},
		"admin":    model.Statuses,
	}, RequireNote: true}
}

func (p Policy) CanChange(role, from, to string) bool {
	if !model.ValidStatus(from) || !model.ValidStatus(to) {
		return false
	}
	allowed, ok := p.AllowedRoles[role]
	if !ok {
		return false
	}
	return contains(allowed, to) && from != to
}

func (p Policy) Check(role string, record model.Record, to, note string) error {
	if !p.CanChange(role, record.Status, to) {
		return fmt.Errorf("role %s cannot move %s to %s", role, record.Status, to)
	}
	if p.RequireNote && strings.TrimSpace(note) == "" {
		return fmt.Errorf("transition note required")
	}
	return nil
}

func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func (p Policy) Roles() []string {
	roles := make([]string, 0, len(p.AllowedRoles))
	for role := range p.AllowedRoles {
		roles = append(roles, role)
	}
	return roles
}

func (p Policy) Describe(role string) string {
	values := p.AllowedRoles[role]
	if len(values) == 0 {
		return "no transitions"
	}
	return strings.Join(values, ",")
}
