package workflow30

import (
	"lawindex/internal/catalog"
	"lawindex/internal/review"
	"lawindex/internal/store"
	"strings"
	"testing"
)

func TestOptionalScoringDisabledMessage(t *testing.T) {
	storage, err := store.Open(t.TempDir() + "/cases.db")
	if err != nil {
		t.Fatal(err)
	}
	defer storage.Close()
	c := catalog.New(storage)
	item, err := c.CreateCase("Trademark review", "client-2", "class-30", "Evidence")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.ValidateCase(item.ID); err != nil {
		t.Fatal(err)
	}
	r := review.New(storage, nil)
	for i := 0; i < 3; i++ {
		_, err = r.RequestScore(item.ID)
	}
	if err == nil || !strings.Contains(err.Error(), "not enabled") {
		t.Fatalf("expected disabled message, got %v", err)
	}
}
