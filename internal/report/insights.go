package report

import (
	"lawindex/internal/model"
	"sort"
)

type Insight struct {
	Label string
	Value float64
	Count int
}

func Insights(cases []model.CaseRecord) []Insight {
	summary := Build(cases)
	result := []Insight{{"completion", CompletionRate(summary), summary.CompletedCount()}}
	for category, count := range summary.ByCategory {
		result = append(result, Insight{Label: category, Value: CategoryShare(summary, category), Count: count})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Label < result[j].Label })
	return result
}
func (s Summary) CompletedCount() int { return s.Approved + s.Rejected + s.Archived }
func AverageScore(cases []model.CaseRecord) float64 {
	total, count := 0, 0
	for _, c := range cases {
		if c.Score > 0 {
			total += c.Score
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return float64(total) / float64(count)
}
func HighestScore(cases []model.CaseRecord) int {
	high := 0
	for _, c := range cases {
		if c.Score > high {
			high = c.Score
		}
	}
	return high
}
