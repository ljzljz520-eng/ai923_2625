package model

import (
	"fmt"
	"strings"
	"time"
)

func FormatCase(c CaseRecord) string {
	return fmt.Sprintf("%s [%s] %s", c.Title, StatusLabel(c.Status), c.Category)
}
func FormatClient(c ClientProfile) string { return fmt.Sprintf("%s <%s>", c.Name, c.Contact) }
func FormatReview(r ReviewRecord) string {
	return fmt.Sprintf("%s by %s", DecisionLabel(r.Decision), r.Reviewer)
}
func FormatArchive(a ArchiveEntry) string { return fmt.Sprintf("%s: %s", a.ArchivedBy, a.Reason) }
func CaseAge(c CaseRecord, now time.Time) time.Duration {
	if c.CreatedAt.IsZero() {
		return 0
	}
	if now.Before(c.CreatedAt) {
		return 0
	}
	return now.Sub(c.CreatedAt)
}
func CompactText(value string, limit int) string {
	value = strings.TrimSpace(value)
	if limit <= 0 {
		return ""
	}
	if len(value) <= limit {
		return value
	}
	return value[:limit]
}
func JoinLabels(values []string) string { return strings.Join(values, ", ") }
