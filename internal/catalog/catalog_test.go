package catalog

import (
	"lawindex/internal/model"
	"lawindex/internal/store"
	"testing"
)

func TestCaseValidation(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/c.db")
	defer s.Close()
	c := New(s)
	if _, err := c.CreateCase("", "x", "y", ""); err == nil {
		t.Fatal("missing title accepted")
	}
	item, err := c.CreateCase("A", "client", "class-30", "summary")
	if err != nil {
		t.Fatal(err)
	}
	if err := c.ValidateCase(item.ID); err != nil {
		t.Fatal(err)
	}
	if err := c.ValidateCase(item.ID); err == nil {
		t.Fatal("submitted case revalidated")
	}
}
func TestFindCasesFilters(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/c.db")
	defer s.Close()
	c := New(s)
	first, _ := c.CreateCase("Alpha", "a", "class-30", "Lease")
	_, _ = c.CreateCase("Beta", "b", "class-45", "Patent")
	_ = c.ValidateCase(first.ID)
	got, err := c.FindCases(model.SearchFilter{Category: "class-30", Status: "submitted", Query: "lease"})
	if err != nil || len(got) != 1 {
		t.Fatalf("got %d, err %v", len(got), err)
	}
}
