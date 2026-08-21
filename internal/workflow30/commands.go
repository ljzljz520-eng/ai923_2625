package workflow30

import (
	"errors"
	"lawindex/internal/model"
)

type Command struct {
	Name   string
	CaseID string
	Actor  string
	Reason string
}

func (w Chain) Execute(command Command) error {
	switch command.Name {
	case "submit":
		return w.Submit(command.CaseID)
	case "archive":
		_, err := w.Close(command.CaseID, command.Actor)
		return err
	case "restore":
		return w.Archive.Restore(command.CaseID)
	default:
		return errors.New("unknown workflow command")
	}
}
func (w Chain) Snapshot(caseID string) (model.CaseRecord, error) {
	return w.Catalog.CaseSummary(caseID)
}
func (w Chain) Ready(caseID string) bool  { return w.Catalog.IsReadyForReview(caseID) }
func (w Chain) Closed(caseID string) bool { return w.Catalog.IsClosed(caseID) }
func CommandNames() []string              { return []string{"submit", "archive", "restore"} }
func ValidCommand(name string) bool {
	for _, candidate := range CommandNames() {
		if name == candidate {
			return true
		}
	}
	return false
}
