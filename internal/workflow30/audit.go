package workflow30

import (
	"fmt"
	"sort"
	"time"
)

type AuditEvent struct {
	ID     string
	CaseID string
	Action string
	Actor  string
	At     time.Time
	Detail string
}
type AuditLog struct{ events []AuditEvent }

func NewAuditLog() *AuditLog { return &AuditLog{events: []AuditEvent{}} }
func (l *AuditLog) Add(event AuditEvent) {
	if event.ID == "" {
		event.ID = fmt.Sprintf("audit-%d", len(l.events)+1)
	}
	l.events = append(l.events, event)
}
func (l *AuditLog) Events(caseID string) []AuditEvent {
	result := []AuditEvent{}
	for _, event := range l.events {
		if caseID == "" || event.CaseID == caseID {
			result = append(result, event)
		}
	}
	sort.Slice(result, func(i, j int) bool { return result[i].At.Before(result[j].At) })
	return result
}
func (l *AuditLog) Count(caseID string) int { return len(l.Events(caseID)) }
func (l *AuditLog) Last(caseID string) (AuditEvent, bool) {
	events := l.Events(caseID)
	if len(events) == 0 {
		return AuditEvent{}, false
	}
	return events[len(events)-1], true
}
func (l *AuditLog) Clear() { l.events = nil }
