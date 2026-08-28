package model

import "fmt"

func ValidateRecord(r Record) error {
	if r.ID == "" {
		return fmt.Errorf("record id required")
	}
	if r.EquipmentID == "" {
		return fmt.Errorf("equipment required")
	}
	if r.Title == "" {
		return fmt.Errorf("title required")
	}
	if !ValidStatus(r.Status) {
		return fmt.Errorf("invalid status")
	}
	return nil
}
func ValidateUser(u User) error {
	if u.ID == "" || u.Name == "" {
		return fmt.Errorf("user identity required")
	}
	if u.Role == "" {
		return fmt.Errorf("role required")
	}
	return nil
}
func ValidateEquipment(e Equipment) error {
	if e.ID == "" || e.Name == "" {
		return fmt.Errorf("equipment identity required")
	}
	return nil
}
