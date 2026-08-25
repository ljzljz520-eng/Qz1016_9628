package model

import "errors"

var errInvalidPlan = errors.New("invalid maintenance plan")

func PriorityWeight(priority string) int {
	switch priority {
	case "urgent":
		return 4
	case "high":
		return 3
	case "normal":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}
func ComparePriority(left, right string) int {
	a, b := PriorityWeight(left), PriorityWeight(right)
	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}
func CanAssign(role, status string) bool {
	if role == "operator" {
		return status == "validated" || status == "processing" || status == "blocked"
	}
	if role == "reviewer" {
		return status == "validated" || status == "resolved"
	}
	return role == "admin"
}
