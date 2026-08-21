package review

import "lawindex/internal/model"

type RuleScorer struct{ Weight int }

func (r RuleScorer) Score(c model.CaseRecord) int {
	score := r.Weight
	if len(c.Summary) > 100 {
		score += 20
	}
	if c.Category == "class-30" {
		score += 10
	}
	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}
func (s *Service) ScoreCase(caseID string) (int, error) { return s.RequestScore(caseID) }
func ScoreLabel(score int) string {
	if score >= 80 {
		return "high"
	}
	if score >= 50 {
		return "medium"
	}
	return "low"
}
func ClampScore(score int) int {
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}
	return score
}
func ScoreRange(score int) bool { return score >= 0 && score <= 100 }
