package workflow30

import (
	"errors"
	"lawindex/internal/model"
)

type Stage string

const (
	StageIntake     Stage = "intake"
	StageSubmission Stage = "submission"
	StageReview     Stage = "review"
	StageArchive    Stage = "archive"
)

func StageFor(status string) Stage {
	switch status {
	case "draft":
		return StageIntake
	case "submitted":
		return StageSubmission
	case "approved", "rejected":
		return StageReview
	case "archived":
		return StageArchive
	default:
		return "unknown"
	}
}
func NextStage(stage Stage) Stage {
	switch stage {
	case StageIntake:
		return StageSubmission
	case StageSubmission:
		return StageReview
	case StageReview:
		return StageArchive
	default:
		return "complete"
	}
}
func CanAdvance(stage Stage) bool {
	return stage == StageIntake || stage == StageSubmission || stage == StageReview
}
func ValidateChain(c model.CaseRecord) error {
	if c.Status == "" {
		return errors.New("case status is required")
	}
	if !model.IsValidStatus(c.Status) {
		return errors.New("unknown case status")
	}
	return nil
}
func ChainDescription() []string    { return []string{"intake", "submission", "review", "archive"} }
func IsTerminal(status string) bool { return status == "archived" }
