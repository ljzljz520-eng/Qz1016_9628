package model

import "testing"

func TestStatusRules(t *testing.T) {
	if !ValidStatus("processing") || NextStatus("new", true) != "validated" {
		t.Fatal("status")
	}
	if StatusRank("closed") < StatusRank("new") {
		t.Fatal("rank")
	}
}
