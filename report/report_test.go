package report

import (
	"labops/model"
	"testing"
)

func TestFormatRecord(t *testing.T) {
	if FormatRecord(model.Record{ID: "r", Status: "new", Title: "T"}) != "r [new] T" {
		t.Fatal("format")
	}
}
