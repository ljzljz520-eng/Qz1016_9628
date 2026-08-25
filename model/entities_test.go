package model

import (
	"testing"
	"time"
)

func TestNewRecord(t *testing.T) {
	r := NewRecord("r1", "e1", "title", time.Now())
	if r.Status != "new" {
		t.Fatal(r.Status)
	}
}
