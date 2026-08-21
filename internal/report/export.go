package report

import (
	"encoding/json"
	"fmt"
	"lawindex/internal/model"
	"sort"
	"strings"
)

func BuildJSON(cases []model.CaseRecord) ([]byte, error) {
	summary := Build(cases)
	return json.Marshal(summary)
}
func BuildCSV(cases []model.CaseRecord) string {
	lines := []string{"id,title,client,category,status,score"}
	sort.Slice(cases, func(i, j int) bool { return cases[i].ID < cases[j].ID })
	for _, item := range cases {
		lines = append(lines, strings.Join([]string{item.ID, quote(item.Title), item.ClientID, item.Category, item.Status, fmt.Sprint(item.Score)}, ","))
	}
	return strings.Join(lines, "\n") + "\n"
}
func quote(value string) string {
	if strings.ContainsAny(value, ",\"\n") {
		return "\"" + strings.ReplaceAll(value, "\"", "\"\"") + "\""
	}
	return value
}
func StatusBreakdown(cases []model.CaseRecord) []string {
	summary := Build(cases)
	result := []string{}
	for _, status := range model.Statuses() {
		result = append(result, fmt.Sprintf("%s:%d", status, statusCount(summary, status)))
	}
	return result
}
func statusCount(summary Summary, status string) int {
	switch status {
	case "draft":
		return summary.Draft
	case "submitted":
		return summary.Submitted
	case "approved":
		return summary.Approved
	case "rejected":
		return summary.Rejected
	case "archived":
		return summary.Archived
	default:
		return 0
	}
}
func CategoryShare(summary Summary, category string) float64 {
	if summary.Total == 0 {
		return 0
	}
	return float64(summary.ByCategory[category]) / float64(summary.Total)
}
