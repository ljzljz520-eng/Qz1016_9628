package model

var Statuses = []string{"new", "validated", "processing", "blocked", "resolved", "closed", "archived", "cancelled"}

func ValidStatus(s string) bool {
	for _, v := range Statuses {
		if s == v {
			return true
		}
	}
	return false
}
func StatusRank(s string) int {
	for i, v := range Statuses {
		if s == v {
			return i
		}
	}
	return -1
}
func NormalizePriority(p string) string {
	switch p {
	case "urgent", "high", "normal", "low":
		return p
	}
	return "normal"
}
func NextStatus(current string, approve bool) string {
	if current == "new" {
		if approve {
			return "validated"
		}
		return "cancelled"
	}
	if current == "validated" {
		return "processing"
	}
	if current == "processing" {
		if approve {
			return "resolved"
		}
		return "blocked"
	}
	if current == "resolved" {
		return "closed"
	}
	if current == "closed" {
		return "archived"
	}
	return current
}
