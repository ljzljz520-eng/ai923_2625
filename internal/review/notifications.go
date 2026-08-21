package review

import (
	"fmt"
	"lawindex/internal/model"
)

type Notification struct {
	CaseID    string
	Recipient string
	Subject   string
	Body      string
}

func BuildNotification(c model.CaseRecord, r model.ReviewRecord) Notification {
	return Notification{CaseID: c.ID, Recipient: r.Reviewer, Subject: "Case " + c.Title, Body: fmt.Sprintf("Decision: %s", model.DecisionLabel(r.Decision))}
}
func (s *Service) NotificationFor(caseID string) (Notification, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return Notification{}, err
	}
	rs, err := s.storage.ListReviews(caseID)
	if err != nil || len(rs) == 0 {
		return Notification{}, fmt.Errorf("review not found")
	}
	return BuildNotification(c, rs[len(rs)-1]), nil
}
func (s *Service) ReviewCount(caseID string) int {
	rs, err := s.storage.ListReviews(caseID)
	if err != nil {
		return 0
	}
	return len(rs)
}
func (s *Service) Approved(caseID string) bool {
	c, err := s.storage.GetCase(caseID)
	return err == nil && c.Status == "approved"
}
func (s *Service) Rejected(caseID string) bool {
	c, err := s.storage.GetCase(caseID)
	return err == nil && c.Status == "rejected"
}
