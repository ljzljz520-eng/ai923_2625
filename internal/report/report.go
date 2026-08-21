package report

import (
	"lawindex/internal/model"
	"sort"
)

type Summary struct {
	Total, Draft, Submitted, Approved, Rejected, Archived int
	ByCategory                                            map[string]int
}

func Build(cases []model.CaseRecord) Summary {
	out := Summary{ByCategory: map[string]int{}}
	for _, c := range cases {
		out.Total++
		out.ByCategory[c.Category]++
		switch c.Status {
		case "draft":
			out.Draft++
		case "submitted":
			out.Submitted++
		case "approved":
			out.Approved++
		case "rejected":
			out.Rejected++
		case "archived":
			out.Archived++
		}
	}
	return out
}
func TopCategories(summary Summary, limit int) []string {
	keys := make([]string, 0, len(summary.ByCategory))
	for k := range summary.ByCategory {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if summary.ByCategory[keys[i]] == summary.ByCategory[keys[j]] {
			return keys[i] < keys[j]
		}
		return summary.ByCategory[keys[i]] > summary.ByCategory[keys[j]]
	})
	if limit < len(keys) {
		return keys[:limit]
	}
	return keys
}
func CompletionRate(summary Summary) float64 {
	if summary.Total == 0 {
		return 0
	}
	return float64(summary.Approved+summary.Rejected+summary.Archived) / float64(summary.Total)
}
