package model

import "strings"

func CanonicalTitle(title string) string {
	return strings.TrimSpace(strings.Join(strings.Fields(title), " "))
}
func IsHighPriority(p string) bool { return p == "urgent" || p == "high" }
func CloneRecord(r Record) Record  { return r }
