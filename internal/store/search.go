package store

import (
	"lawindex/internal/model"
	"sort"
	"strings"
)

type CasePage struct {
	Items  []model.CaseRecord
	Offset int
	Limit  int
	Total  int
}
type CaseQuery struct {
	Text       string
	Category   string
	Status     string
	Offset     int
	Limit      int
	Descending bool
}

func (s *Store) QueryCases(query CaseQuery) (CasePage, error) {
	all, err := s.ListCases()
	if err != nil {
		return CasePage{}, err
	}
	filtered := make([]model.CaseRecord, 0, len(all))
	for _, item := range all {
		if !matchesCase(item, query) {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Slice(filtered, func(i, j int) bool {
		if query.Descending {
			return filtered[i].UpdatedAt.After(filtered[j].UpdatedAt)
		}
		return filtered[i].UpdatedAt.Before(filtered[j].UpdatedAt)
	})
	page := paginate(filtered, query.Offset, query.Limit)
	return CasePage{Items: page, Offset: normalizeOffset(query.Offset), Limit: normalizeLimit(query.Limit), Total: len(filtered)}, nil
}
func matchesCase(item model.CaseRecord, query CaseQuery) bool {
	if query.Category != "" && item.Category != query.Category {
		return false
	}
	if query.Status != "" && item.Status != query.Status {
		return false
	}
	if query.Text != "" {
		haystack := strings.ToLower(item.Title + " " + item.Summary + " " + item.ClientID)
		if !strings.Contains(haystack, strings.ToLower(query.Text)) {
			return false
		}
	}
	return true
}
func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 500 {
		return 500
	}
	return limit
}
func paginate(items []model.CaseRecord, offset, limit int) []model.CaseRecord {
	offset = normalizeOffset(offset)
	limit = normalizeLimit(limit)
	if offset >= len(items) {
		return []model.CaseRecord{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return append([]model.CaseRecord(nil), items[offset:end]...)
}
func CountByStatus(items []model.CaseRecord) map[string]int {
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Status]++
	}
	return counts
}
func CountByCategory(items []model.CaseRecord) map[string]int {
	counts := map[string]int{}
	for _, item := range items {
		counts[item.Category]++
	}
	return counts
}
func (s *Store) CaseExists(id string) bool { _, err := s.GetCase(id); return err == nil }
func (s *Store) SaveCaseIfMissing(item model.CaseRecord) error {
	if s.CaseExists(item.ID) {
		return nil
	}
	return s.PutCase(item)
}
