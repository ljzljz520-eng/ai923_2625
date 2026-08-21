package workflow30

import (
	"lawindex/internal/model"
)

type Metrics struct {
	Total        int
	Active       int
	Completed    int
	HighPriority int
	Categories   map[string]int
}

func ComputeMetrics(cases []model.CaseRecord) Metrics {
	result := Metrics{Categories: map[string]int{}}
	for _, item := range cases {
		result.Total++
		result.Categories[item.Category]++
		if item.IsOpen() {
			result.Active++
		}
		if item.IsReviewed() || item.Status == "archived" {
			result.Completed++
		}
		if item.Score >= 80 {
			result.HighPriority++
		}
	}
	return result
}
func (m Metrics) CompletionPercent() float64 {
	if m.Total == 0 {
		return 0
	}
	return float64(m.Completed*100) / float64(m.Total)
}
func (m Metrics) CategoryCount(category string) int { return m.Categories[category] }
func (m Metrics) IsBusy() bool                      { return m.Active >= 10 }
func (m Metrics) Labels() []string                  { return []string{"total", "active", "completed", "high-priority"} }
