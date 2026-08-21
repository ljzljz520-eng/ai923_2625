package archive

import (
	"errors"
	"lawindex/internal/model"
	"strings"
	"time"
)

type RetentionRule struct {
	Days     int
	Category string
}

func (r RetentionRule) Valid() bool { return r.Days > 0 }
func ShouldRetain(entry model.ArchiveEntry, rule RetentionRule, now time.Time) bool {
	if !rule.Valid() || entry.ArchivedAt.IsZero() {
		return false
	}
	if rule.Category != "" && !strings.Contains(entry.Reason, rule.Category) {
		return false
	}
	return now.Sub(entry.ArchivedAt) < time.Duration(rule.Days)*24*time.Hour
}
func ValidateRetention(rule RetentionRule) error {
	if rule.Days < 1 {
		return errors.New("retention days must be positive")
	}
	return nil
}
func (s *Service) ArchiveAge(entry model.ArchiveEntry, now time.Time) time.Duration {
	if now.Before(entry.ArchivedAt) {
		return 0
	}
	return now.Sub(entry.ArchivedAt)
}
func NormalizeActor(actor string) string { return strings.TrimSpace(actor) }
