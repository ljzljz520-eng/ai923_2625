package workflow30

import "lawindex/internal/model"

type Plan struct {
	CaseID  string
	Steps   []string
	Current int
}

func NewPlan(c model.CaseRecord) Plan {
	current := 0
	switch StageFor(c.Status) {
	case StageSubmission:
		current = 1
	case StageReview:
		current = 2
	case StageArchive:
		current = 3
	}
	return Plan{CaseID: c.ID, Steps: ChainDescription(), Current: current}
}
func (p Plan) Complete() bool { return p.Current >= len(p.Steps) }
func (p *Plan) Advance() bool {
	if p.Complete() {
		return false
	}
	p.Current++
	return true
}
func (p Plan) CurrentStep() string {
	if p.Current < 0 || p.Current >= len(p.Steps) {
		return "complete"
	}
	return p.Steps[p.Current]
}
func (p Plan) Remaining() int {
	if p.Complete() {
		return 0
	}
	return len(p.Steps) - p.Current
}
