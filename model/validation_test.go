package model

import (
	"testing"
	"time"
)

func TestValidation(t *testing.T) {
	if ValidateRecord(Record{}) == nil {
		t.Fatal("expected error")
	}
	if ValidateRecord(NewRecord("r", "e", "t", time.Now())) != nil {
		t.Fatal("valid")
	}
}
