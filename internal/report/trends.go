package report

import (
	"lawindex/internal/model"
	"sort"
	"time"
)

type PeriodCount struct {
	Date  string
	Count int
}

func DailyCounts(cases []model.CaseRecord, location *time.Location) []PeriodCount {
	if location == nil {
		location = time.UTC
	}
	counts := map[string]int{}
	for _, c := range cases {
		key := c.CreatedAt.In(location).Format("2006-01-02")
		counts[key]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]PeriodCount, 0, len(keys))
	for _, key := range keys {
		result = append(result, PeriodCount{key, counts[key]})
	}
	return result
}
func GrowthRate(previous, current int) float64 {
	if previous <= 0 {
		if current > 0 {
			return 1
		}
		return 0
	}
	return float64(current-previous) / float64(previous)
}
func MedianScore(cases []model.CaseRecord) float64 {
	scores := []int{}
	for _, c := range cases {
		if c.Score > 0 {
			scores = append(scores, c.Score)
		}
	}
	if len(scores) == 0 {
		return 0
	}
	sort.Ints(scores)
	mid := len(scores) / 2
	if len(scores)%2 == 1 {
		return float64(scores[mid])
	}
	return float64(scores[mid-1]+scores[mid]) / 2
}
func StatusAt(cases []model.CaseRecord, status string) []model.CaseRecord {
	result := []model.CaseRecord{}
	for _, c := range cases {
		if c.Status == status {
			result = append(result, c)
		}
	}
	return result
}
