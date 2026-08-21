package workflow30

import (
	"errors"
	"lawindex/internal/archive"
	"lawindex/internal/catalog"
	"lawindex/internal/model"
	"lawindex/internal/review"
)

type Chain struct {
	Catalog *catalog.Service
	Review  *review.Service
	Archive *archive.Service
}

func NewChain(c *catalog.Service, r *review.Service, a *archive.Service) Chain {
	return Chain{Catalog: c, Review: r, Archive: a}
}
func (w Chain) Intake(title, client, category, summary string) (model.CaseRecord, error) {
	return w.Catalog.CreateCase(title, client, category, summary)
}
func (w Chain) Submit(caseID string) error { return w.Catalog.ValidateCase(caseID) }
func (w Chain) Decide(caseID, reviewer, decision string) (model.ReviewRecord, error) {
	if reviewer == "" {
		return model.ReviewRecord{}, errors.New("reviewer is required")
	}
	return w.Review.SubmitReview(caseID, reviewer, decision, "")
}
func (w Chain) Close(caseID, actor string) (model.ArchiveEntry, error) {
	return w.Archive.Archive(caseID, actor, "completed matter")
}
