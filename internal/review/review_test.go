package review

import (
	"lawindex/internal/catalog"
	"lawindex/internal/model"
	"lawindex/internal/store"
	"testing"
)

type fixedScorer struct{}

func (fixedScorer) Score(model.CaseRecord) int { return 88 }
func TestReviewHistory(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/r.db")
	defer s.Close()
	c := catalog.New(s)
	item, _ := c.CreateCase("Matter", "c", "class-30", "s")
	_ = c.ValidateCase(item.ID)
	r := New(s, fixedScorer{})
	rec, err := r.SubmitReview(item.ID, "lawyer", "approved", "ready")
	if err != nil || rec.Score != 88 {
		t.Fatalf("review: %+v %v", rec, err)
	}
	history, err := r.ReviewHistory(item.ID)
	if err != nil || len(history) != 1 {
		t.Fatalf("history: %d %v", len(history), err)
	}
}
