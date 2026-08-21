package report

import (
	"lawindex/internal/model"
	"testing"
)

func TestReportCompletionRate(t *testing.T) {
	summary := Build([]model.CaseRecord{{Status: "draft", Category: "a"}, {Status: "approved", Category: "a"}, {Status: "archived", Category: "b"}})
	if summary.Total != 3 || CompletionRate(summary) != 2.0/3.0 {
		t.Fatalf("summary: %+v", summary)
	}
	top := TopCategories(summary, 1)
	if len(top) != 1 || top[0] != "a" {
		t.Fatalf("top: %v", top)
	}
}
