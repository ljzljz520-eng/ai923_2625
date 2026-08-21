package catalog

import (
	"errors"
	"lawindex/internal/model"
	"strings"
)

func (s *Service) ReopenCase(id string) error {
	item, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if item.Status != "rejected" {
		return errors.New("only rejected cases can be reopened")
	}
	item.Status = "draft"
	item.UpdatedAt = s.now()
	return s.storage.PutCase(item)
}
func (s *Service) ChangeCategory(id, category string) error {
	category = strings.TrimSpace(category)
	if category == "" {
		return errors.New("category is required")
	}
	item, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if item.Status == "archived" {
		return errors.New("archived case cannot change category")
	}
	item.Category = category
	item.UpdatedAt = s.now()
	return s.storage.PutCase(item)
}
func (s *Service) AssignClient(id, clientID string) error {
	if strings.TrimSpace(clientID) == "" {
		return errors.New("client is required")
	}
	if _, err := s.storage.GetClient(clientID); err != nil {
		return err
	}
	item, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if item.Status != "draft" {
		return errors.New("only draft cases can change client")
	}
	item.ClientID = clientID
	item.UpdatedAt = s.now()
	return s.storage.PutCase(item)
}
func (s *Service) CaseSummary(id string) (model.CaseRecord, error) { return s.storage.GetCase(id) }
func (s *Service) IsReadyForReview(id string) bool {
	item, err := s.storage.GetCase(id)
	return err == nil && item.Status == "submitted"
}
func (s *Service) IsClosed(id string) bool {
	item, err := s.storage.GetCase(id)
	return err == nil && item.Status == "archived"
}
