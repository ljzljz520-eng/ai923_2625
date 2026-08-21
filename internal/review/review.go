package review

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"lawindex/internal/model"
	"lawindex/internal/store"
)

type Scorer interface{ Score(model.CaseRecord) int }
type Service struct {
	storage *store.Store
	scorer  Scorer
	now     func() time.Time
}

func New(storage *store.Store, scorer Scorer) *Service {
	return &Service{storage: storage, scorer: scorer, now: time.Now}
}
func (s *Service) SetScorer(scorer Scorer) { s.scorer = scorer }
func (s *Service) RequestScore(caseID string) (int, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return 0, err
	}
	if s.scorer == nil {
		return 0, errors.New("optional scorer is not enabled")
	}
	return s.scorer.Score(c), nil
}
func (s *Service) SubmitReview(caseID, reviewer, decision, comment string) (model.ReviewRecord, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return model.ReviewRecord{}, err
	}
	if c.Status != "submitted" {
		return model.ReviewRecord{}, errors.New("case is not submitted")
	}
	decision = strings.ToLower(strings.TrimSpace(decision))
	if decision != "approved" && decision != "rejected" {
		return model.ReviewRecord{}, errors.New("decision must be approved or rejected")
	}
	score, scoreErr := s.RequestScore(caseID)
	if scoreErr != nil {
		score = 0
	}
	c.Status = decision
	c.Score = score
	c.UpdatedAt = s.now()
	if err := s.storage.PutCase(c); err != nil {
		return model.ReviewRecord{}, err
	}
	r := model.ReviewRecord{ID: fmt.Sprintf("review-%d", s.now().UnixNano()), CaseID: caseID, Reviewer: strings.TrimSpace(reviewer), Decision: decision, Comment: strings.TrimSpace(comment), Score: score, CreatedAt: s.now()}
	return r, s.storage.PutReview(r)
}
func (s *Service) ReviewHistory(caseID string) ([]model.ReviewRecord, error) {
	return s.storage.ListReviews(caseID)
}
func (s *Service) IsHighPriority(score int) bool { return score >= 80 }
