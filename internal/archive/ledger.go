package archive

import (
	"errors"
	"lawindex/internal/model"
	"sort"
	"strings"
)

func (s *Service) ListArchives() ([]model.ArchiveEntry, error) {
	cases, err := s.storage.ListCases()
	if err != nil {
		return nil, err
	}
	result := []model.ArchiveEntry{}
	for _, item := range cases {
		if item.Status != "archived" {
			continue
		}
		archives, listErr := s.storage.ListReviews(item.ID)
		if listErr != nil {
			return nil, listErr
		}
		_ = archives
	}
	all := []model.ArchiveEntry{}
	for _, item := range cases {
		if item.Status == "archived" {
			entry, getErr := s.storage.GetArchive("archive-" + item.ID)
			if getErr == nil {
				all = append(all, entry)
			}
		}
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ArchivedAt.Before(all[j].ArchivedAt) })
	return resultOrAll(result, all), nil
}
func resultOrAll(empty, all []model.ArchiveEntry) []model.ArchiveEntry {
	if len(all) > 0 {
		return all
	}
	return empty
}
func (s *Service) ArchiveReason(id string) (string, error) {
	entry, err := s.storage.GetArchive(id)
	if err != nil {
		return "", err
	}
	return entry.Reason, nil
}
func ValidateArchiveRequest(actor, reason string) error {
	if strings.TrimSpace(actor) == "" {
		return errors.New("actor is required")
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("reason is required")
	}
	return nil
}
func NormalizeReason(reason string) string { return strings.TrimSpace(reason) }
