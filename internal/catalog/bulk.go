package catalog

import (
	"errors"
	"fmt"
	"lawindex/internal/model"
	"strings"
)

type IntakeRow struct {
	Title    string
	ClientID string
	Category string
	Summary  string
}
type ImportReport struct {
	Accepted int
	Rejected int
	IDs      []string
	Errors   []string
}

func (s *Service) ImportRows(rows []IntakeRow) ImportReport {
	report := ImportReport{IDs: []string{}, Errors: []string{}}
	for index, row := range rows {
		item, err := s.CreateCase(row.Title, row.ClientID, row.Category, row.Summary)
		if err != nil {
			report.Rejected++
			report.Errors = append(report.Errors, fmt.Sprintf("row %d: %v", index+1, err))
			continue
		}
		report.Accepted++
		report.IDs = append(report.IDs, item.ID)
	}
	return report
}
func (s *Service) ValidateRows(rows []IntakeRow) []error {
	failures := []error{}
	for index, row := range rows {
		if strings.TrimSpace(row.Title) == "" {
			failures = append(failures, fmt.Errorf("row %d title is required", index+1))
		}
		if strings.TrimSpace(row.ClientID) == "" {
			failures = append(failures, fmt.Errorf("row %d client is required", index+1))
		}
		if strings.TrimSpace(row.Category) == "" {
			failures = append(failures, fmt.Errorf("row %d category is required", index+1))
		}
	}
	return failures
}
func ParseRow(parts []string) (IntakeRow, error) {
	if len(parts) < 3 {
		return IntakeRow{}, errors.New("at least title, client, and category are required")
	}
	row := IntakeRow{Title: strings.TrimSpace(parts[0]), ClientID: strings.TrimSpace(parts[1]), Category: strings.TrimSpace(parts[2])}
	if len(parts) > 3 {
		row.Summary = strings.TrimSpace(parts[3])
	}
	if err := (&Service{}).ValidateRows([]IntakeRow{row}); len(err) > 0 {
		return IntakeRow{}, err[0]
	}
	return row, nil
}
func CaseToRow(item model.CaseRecord) IntakeRow {
	return IntakeRow{Title: item.Title, ClientID: item.ClientID, Category: item.Category, Summary: item.Summary}
}
func RowKey(row IntakeRow) string {
	return strings.ToLower(strings.Join([]string{row.Title, row.ClientID, row.Category}, "|"))
}
func DeduplicateRows(rows []IntakeRow) []IntakeRow {
	seen := map[string]bool{}
	result := []IntakeRow{}
	for _, row := range rows {
		key := RowKey(row)
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, row)
	}
	return result
}
func MergeSummary(existing, incoming string) string {
	existing = strings.TrimSpace(existing)
	incoming = strings.TrimSpace(incoming)
	if existing == "" {
		return incoming
	}
	if incoming == "" {
		return existing
	}
	return existing + "\n" + incoming
}
