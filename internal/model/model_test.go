package model

import (
	"testing"
	"time"
)

func TestCaseStatePredicates(t *testing.T) {
	c := CaseRecord{Status: "approved"}
	if !c.IsOpen() || !c.IsReviewed() || !c.CanArchive() {
		t.Fatal("approved state predicates failed")
	}
	c.Status = "archived"
	if c.IsOpen() || c.CanArchive() {
		t.Fatal("archived state predicates failed")
	}
	c.Touch(time.Unix(10, 0))
	if c.UpdatedAt.IsZero() {
		t.Fatal("touch failed")
	}
}
