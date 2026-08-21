package archive

import (
	"lawindex/internal/catalog"
	"lawindex/internal/model"
	"lawindex/internal/review"
	"lawindex/internal/store"
	"testing"
)

func TestWorkflowArchiveLifecycle(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/a.db")
	defer s.Close()
	c := catalog.New(s)
	item, _ := c.CreateCase("Matter", "c", "class-30", "s")
	if _, err := New(s).Archive(item.ID, "admin", "done"); err == nil {
		t.Fatal("draft archived")
	}
	_ = c.ValidateCase(item.ID)
	_, _ = review.New(s, fixed{}).SubmitReview(item.ID, "lawyer", "approved", "")
	a := New(s)
	entry, err := a.Archive(item.ID, "admin", "done")
	if err != nil || entry.CaseID != item.ID {
		t.Fatalf("archive: %+v %v", entry, err)
	}
	if err := a.Restore(item.ID); err != nil {
		t.Fatal(err)
	}
}

type fixed struct{}

func (fixed) Score(model.CaseRecord) int { return 70 }
