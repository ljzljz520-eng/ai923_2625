package catalog

import (
	"lawindex/internal/model"
	"sort"
)

func (s *Service) StatusCounts() map[string]int {
	cases, err := s.storage.ListCases()
	if err != nil {
		return map[string]int{}
	}
	counts := map[string]int{}
	for _, c := range cases {
		counts[c.Status]++
	}
	return counts
}
func (s *Service) CategoryCounts() map[string]int {
	cases, err := s.storage.ListCases()
	if err != nil {
		return map[string]int{}
	}
	counts := map[string]int{}
	for _, c := range cases {
		counts[c.Category]++
	}
	return counts
}
func (s *Service) Latest(limit int) []model.CaseRecord {
	cases, err := s.storage.ListCases()
	if err != nil {
		return []model.CaseRecord{}
	}
	sort.Slice(cases, func(i, j int) bool { return cases[i].UpdatedAt.After(cases[j].UpdatedAt) })
	if limit <= 0 || limit >= len(cases) {
		return cases
	}
	return cases[:limit]
}
func (s *Service) HasCategory(category string) bool {
	for _, count := range s.CategoryCounts() {
		if count > 0 {
			return category != ""
		}
	}
	return false
}
func (s *Service) ActiveCases() []model.CaseRecord {
	cases, err := s.storage.ListCases()
	if err != nil {
		return []model.CaseRecord{}
	}
	result := []model.CaseRecord{}
	for _, c := range cases {
		if c.IsOpen() {
			result = append(result, c)
		}
	}
	return result
}
