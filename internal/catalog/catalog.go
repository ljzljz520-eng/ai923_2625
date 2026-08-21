package catalog

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"lawindex/internal/model"
	"lawindex/internal/store"
)

type Service struct {
	storage *store.Store
	now     func() time.Time
}

func New(storage *store.Store) *Service { return &Service{storage: storage, now: time.Now} }

func (s *Service) CreateCase(title, clientID, category, summary string) (model.CaseRecord, error) {
	if strings.TrimSpace(title) == "" {
		return model.CaseRecord{}, errors.New("title is required")
	}
	if strings.TrimSpace(clientID) == "" {
		return model.CaseRecord{}, errors.New("client is required")
	}
	if strings.TrimSpace(category) == "" {
		return model.CaseRecord{}, errors.New("category is required")
	}
	now := s.now()
	c := model.CaseRecord{ID: fmt.Sprintf("case-%d", now.UnixNano()), Title: strings.TrimSpace(title), ClientID: clientID, Category: category, Summary: strings.TrimSpace(summary), Status: "draft", CreatedAt: now, UpdatedAt: now}
	return c, s.storage.PutCase(c)
}
func (s *Service) RegisterClient(name, contact, notes string) (model.ClientProfile, error) {
	if strings.TrimSpace(name) == "" {
		return model.ClientProfile{}, errors.New("name is required")
	}
	c := model.ClientProfile{ID: fmt.Sprintf("client-%d", s.now().UnixNano()), Name: strings.TrimSpace(name), Contact: strings.TrimSpace(contact), Notes: strings.TrimSpace(notes)}
	return c, s.storage.PutClient(c)
}
func (s *Service) FindCases(filter model.SearchFilter) ([]model.CaseRecord, error) {
	cases, err := s.storage.ListCases()
	if err != nil {
		return nil, err
	}
	var result []model.CaseRecord
	for _, c := range cases {
		if filter.Category != "" && c.Category != filter.Category {
			continue
		}
		if filter.Status != "" && c.Status != filter.Status {
			continue
		}
		if filter.Query != "" && !strings.Contains(strings.ToLower(c.Title+" "+c.Summary), strings.ToLower(filter.Query)) {
			continue
		}
		result = append(result, c)
	}
	return result, nil
}
func (s *Service) UpdateSummary(id, summary string) error {
	c, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if c.Status == "archived" {
		return errors.New("archived case cannot be edited")
	}
	c.Summary = strings.TrimSpace(summary)
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
func (s *Service) ValidateCase(id string) error {
	c, err := s.storage.GetCase(id)
	if err != nil {
		return err
	}
	if c.Title == "" || c.ClientID == "" || c.Category == "" {
		return errors.New("case is incomplete")
	}
	if c.Status != "draft" {
		return errors.New("case is not awaiting review")
	}
	c.Status = "submitted"
	c.UpdatedAt = s.now()
	return s.storage.PutCase(c)
}
