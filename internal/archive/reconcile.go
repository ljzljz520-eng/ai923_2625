package archive

import (
	"errors"
	"lawindex/internal/model"
	"strings"
)

func ValidateArchiveEntry(entry model.ArchiveEntry) error {
	if entry.ID == "" || entry.CaseID == "" {
		return errors.New("archive identity is required")
	}
	if strings.TrimSpace(entry.Reason) == "" {
		return errors.New("archive reason is required")
	}
	return nil
}
func (s *Service) Reconcile(caseID string) (bool, error) {
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return false, err
	}
	if c.Status != "archived" {
		return false, nil
	}
	return true, nil
}
func (s *Service) CanRestore(caseID string) bool {
	c, err := s.storage.GetCase(caseID)
	return err == nil && c.Status == "archived"
}
func (s *Service) RestoreTo(caseID, status string) error {
	if status != "approved" && status != "rejected" {
		return errors.New("restore status must be reviewed")
	}
	c, err := s.storage.GetCase(caseID)
	if err != nil {
		return err
	}
	if c.Status != "archived" {
		return errors.New("case is not archived")
	}
	c.Status = status
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
func ArchiveStatuses() []string { return []string{"approved", "rejected", "archived"} }
