package catalog

import (
	"lawindex/internal/model"
	"sort"
	"strings"
	"time"
)

type TimelineEvent struct {
	At     time.Time
	Label  string
	Detail string
}

func BuildTimeline(c model.CaseRecord, reviews []model.ReviewRecord, archive *model.ArchiveEntry) []TimelineEvent {
	events := []TimelineEvent{}
	if !c.CreatedAt.IsZero() {
		events = append(events, TimelineEvent{c.CreatedAt, "created", c.Title})
	}
	if !c.UpdatedAt.IsZero() && c.UpdatedAt != c.CreatedAt {
		events = append(events, TimelineEvent{c.UpdatedAt, "updated", model.StatusLabel(c.Status)})
	}
	for _, r := range reviews {
		events = append(events, TimelineEvent{r.CreatedAt, "reviewed", r.Decision + " by " + r.Reviewer})
	}
	if archive != nil {
		events = append(events, TimelineEvent{archive.ArchivedAt, "archived", archive.Reason})
	}
	sort.Slice(events, func(i, j int) bool { return events[i].At.Before(events[j].At) })
	return events
}
func TimelineLabels(events []TimelineEvent) []string {
	labels := make([]string, 0, len(events))
	for _, e := range events {
		labels = append(labels, e.Label)
	}
	return labels
}
func FilterTimeline(events []TimelineEvent, from, to time.Time) []TimelineEvent {
	result := []TimelineEvent{}
	for _, e := range events {
		if !from.IsZero() && e.At.Before(from) {
			continue
		}
		if !to.IsZero() && e.At.After(to) {
			continue
		}
		result = append(result, e)
	}
	return result
}
func TimelineText(events []TimelineEvent) string {
	parts := make([]string, 0, len(events))
	for _, e := range events {
		parts = append(parts, e.Label+": "+strings.TrimSpace(e.Detail))
	}
	return strings.Join(parts, " | ")
}
func LatestEvent(events []TimelineEvent) (TimelineEvent, bool) {
	if len(events) == 0 {
		return TimelineEvent{}, false
	}
	return events[len(events)-1], true
}
func (s *Service) CaseTimeline(id string) ([]TimelineEvent, error) {
	c, err := s.storage.GetCase(id)
	if err != nil {
		return nil, err
	}
	reviews, err := s.storage.ListReviews(id)
	if err != nil {
		return nil, err
	}
	return BuildTimeline(c, reviews, nil), nil
}
