package review

import (
	"errors"
	"lawindex/internal/model"
	"strings"
)

type Policy struct {
	MinimumComment   int
	RequireReviewer  bool
	AllowedDecisions map[string]bool
}

func DefaultPolicy() Policy {
	return Policy{MinimumComment: 4, RequireReviewer: true, AllowedDecisions: map[string]bool{"approved": true, "rejected": true}}
}
func (p Policy) Validate(r model.ReviewRecord) error {
	if p.RequireReviewer && strings.TrimSpace(r.Reviewer) == "" {
		return errors.New("reviewer is required")
	}
	if !p.AllowedDecisions[r.Decision] {
		return errors.New("decision is not allowed")
	}
	if len(strings.TrimSpace(r.Comment)) < p.MinimumComment {
		return errors.New("review comment is too short")
	}
	return nil
}
func (s *Service) DecideWithPolicy(caseID, reviewer, decision, comment string, policy Policy) (model.ReviewRecord, error) {
	item, err := s.storage.GetCase(caseID)
	if err != nil {
		return model.ReviewRecord{}, err
	}
	candidate := model.ReviewRecord{CaseID: caseID, Reviewer: reviewer, Decision: strings.ToLower(strings.TrimSpace(decision)), Comment: comment}
	if err := policy.Validate(candidate); err != nil {
		return model.ReviewRecord{}, err
	}
	return s.SubmitReview(item.ID, reviewer, candidate.Decision, comment)
}
func (s *Service) DecisionAllowed(decision string) bool {
	normalized := strings.ToLower(strings.TrimSpace(decision))
	return normalized == "approved" || normalized == "rejected"
}
func (s *Service) ScoreAvailable() bool { return s.scorer != nil }
func (s *Service) ScoringStatus() string {
	if s.ScoreAvailable() {
		return "enabled"
	}
	return "disabled"
}
func (s *Service) ExplainScoring() string {
	if s.ScoreAvailable() {
		return "scoring is available"
	}
	return "optional scorer is not enabled"
}
