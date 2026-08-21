package model

import (
	"fmt"
	"strings"
)

var validStatuses = []string{"draft", "submitted", "approved", "rejected", "archived"}

func Statuses() []string { return append([]string(nil), validStatuses...) }
func IsValidStatus(status string) bool {
	for _, value := range validStatuses {
		if status == value {
			return true
		}
	}
	return false
}
func NormalizeStatus(status string) string { return strings.ToLower(strings.TrimSpace(status)) }
func ValidateCase(c CaseRecord) error {
	if strings.TrimSpace(c.ID) == "" {
		return fmt.Errorf("case id is required")
	}
	if strings.TrimSpace(c.Title) == "" {
		return fmt.Errorf("case title is required")
	}
	if strings.TrimSpace(c.ClientID) == "" {
		return fmt.Errorf("client id is required")
	}
	if strings.TrimSpace(c.Category) == "" {
		return fmt.Errorf("case category is required")
	}
	if !IsValidStatus(c.Status) {
		return fmt.Errorf("unknown case status: %s", c.Status)
	}
	return nil
}
func ValidateClient(c ClientProfile) error {
	if c.ID == "" {
		return fmt.Errorf("client id is required")
	}
	if strings.TrimSpace(c.Name) == "" {
		return fmt.Errorf("client name is required")
	}
	return nil
}
func ValidateReview(r ReviewRecord) error {
	if r.ID == "" || r.CaseID == "" {
		return fmt.Errorf("review identity is required")
	}
	if r.Decision != "approved" && r.Decision != "rejected" {
		return fmt.Errorf("invalid review decision")
	}
	if r.Reviewer == "" {
		return fmt.Errorf("reviewer is required")
	}
	return nil
}
func ValidateArchive(a ArchiveEntry) error {
	if a.ID == "" || a.CaseID == "" {
		return fmt.Errorf("archive identity is required")
	}
	if a.ArchivedBy == "" {
		return fmt.Errorf("archived by is required")
	}
	return nil
}
func StatusLabel(status string) string {
	switch NormalizeStatus(status) {
	case "draft":
		return "Draft"
	case "submitted":
		return "Awaiting review"
	case "approved":
		return "Approved"
	case "rejected":
		return "Rejected"
	case "archived":
		return "Archived"
	default:
		return "Unknown"
	}
}
func DecisionLabel(decision string) string {
	if NormalizeStatus(decision) == "approved" {
		return "Approved"
	}
	if NormalizeStatus(decision) == "rejected" {
		return "Rejected"
	}
	return "Pending"
}
