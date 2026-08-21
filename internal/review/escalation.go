package review

import (
	"fmt"
	"lawindex/internal/model"
	"strings"
)

type Escalation struct {
	CaseID string
	Reason string
	To     string
	Urgent bool
}

func BuildEscalation(c model.CaseRecord, score int) Escalation {
	urgent := score >= 80 || strings.Contains(strings.ToLower(c.Summary), "urgent")
	reason := "standard review"
	if urgent {
		reason = "priority review"
	}
	return Escalation{CaseID: c.ID, Reason: reason, To: "senior-attorney", Urgent: urgent}
}
func EscalationText(e Escalation) string {
	level := "normal"
	if e.Urgent {
		level = "urgent"
	}
	return fmt.Sprintf("%s escalation to %s: %s", level, e.To, e.Reason)
}
func (s *Service) Escalate(caseID string) (Escalation, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return Escalation{}, err
	}
	score := c.Score
	if score == 0 && s.scorer != nil {
		score, _ = s.RequestScore(caseID)
	}
	return BuildEscalation(c, score), nil
}
func (s *Service) NeedsEscalation(caseID string) bool {
	e, err := s.Escalate(caseID)
	return err == nil && e.Urgent
}
func EscalationRecipients(e Escalation) []string {
	if e.Urgent {
		return []string{e.To, "managing-partner"}
	}
	return []string{e.To}
}
