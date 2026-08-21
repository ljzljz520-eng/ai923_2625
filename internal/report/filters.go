package report

import (
	"lawindex/internal/model"
	"sort"
	"strings"
)

func FilterByText(cases []model.CaseRecord, text string) []model.CaseRecord {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return append([]model.CaseRecord(nil), cases...)
	}
	result := []model.CaseRecord{}
	for _, item := range cases {
		value := strings.ToLower(item.Title + " " + item.Summary + " " + item.Category)
		if strings.Contains(value, text) {
			result = append(result, item)
		}
	}
	return result
}
func FilterByScore(cases []model.CaseRecord, minimum, maximum int) []model.CaseRecord {
	result := []model.CaseRecord{}
	for _, item := range cases {
		if item.Score < minimum {
			continue
		}
		if maximum > 0 && item.Score > maximum {
			continue
		}
		result = append(result, item)
	}
	return result
}
func SortByScore(cases []model.CaseRecord, descending bool) []model.CaseRecord {
	result := append([]model.CaseRecord(nil), cases...)
	sort.SliceStable(result, func(i, j int) bool {
		if descending {
			return result[i].Score > result[j].Score
		}
		return result[i].Score < result[j].Score
	})
	return result
}
func SortByTitle(cases []model.CaseRecord) []model.CaseRecord {
	result := append([]model.CaseRecord(nil), cases...)
	sort.SliceStable(result, func(i, j int) bool { return strings.ToLower(result[i].Title) < strings.ToLower(result[j].Title) })
	return result
}
func UniqueCategories(cases []model.CaseRecord) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, item := range cases {
		if seen[item.Category] {
			continue
		}
		seen[item.Category] = true
		result = append(result, item.Category)
	}
	sort.Strings(result)
	return result
}
func HasStatus(cases []model.CaseRecord, status string) bool {
	for _, item := range cases {
		if item.Status == status {
			return true
		}
	}
	return false
}
func CountHighPriority(cases []model.CaseRecord, threshold int) int {
	count := 0
	for _, item := range cases {
		if item.Score >= threshold {
			count++
		}
	}
	return count
}
func Titles(cases []model.CaseRecord) []string {
	values := make([]string, 0, len(cases))
	for _, item := range cases {
		values = append(values, item.Title)
	}
	return values
}
