package catalog

import (
	"errors"
	"lawindex/internal/model"
	"strings"
)

type Assignment struct {
	CaseID string
	Staff  string
	Role   string
}

func ValidateAssignment(a Assignment) error {
	if a.CaseID == "" {
		return errors.New("case is required")
	}
	if strings.TrimSpace(a.Staff) == "" {
		return errors.New("staff is required")
	}
	if strings.TrimSpace(a.Role) == "" {
		return errors.New("role is required")
	}
	return nil
}
func (s *Service) AssignStaff(id, staff, role string) error {
	if err := ValidateAssignment(Assignment{id, staff, role}); err != nil {
		return err
	}
	c, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if c.Status == "archived" {
		return errors.New("archived case cannot be assigned")
	}
	c.Summary = MergeSummary(c.Summary, "Assigned "+strings.TrimSpace(staff)+" as "+strings.TrimSpace(role))
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
func (s *Service) UnassignStaff(id, staff string) error {
	c, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if !strings.Contains(c.Summary, staff) {
		return errors.New("staff assignment not found")
	}
	c.Summary = strings.ReplaceAll(c.Summary, "Assigned "+staff, "Unassigned "+staff)
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
func AssignmentLabel(a Assignment) string {
	return strings.TrimSpace(a.Staff) + " / " + strings.TrimSpace(a.Role)
}
func SameAssignment(a, b Assignment) bool {
	return a.CaseID == b.CaseID && strings.EqualFold(a.Staff, b.Staff) && strings.EqualFold(a.Role, b.Role)
}
func CaseNeedsAssignment(c model.CaseRecord) bool {
	return c.Status == "submitted" && !strings.Contains(strings.ToLower(c.Summary), "assigned")
}
