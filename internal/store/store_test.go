package store

import (
	"lawindex/internal/model"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := t.TempDir() + "/persist.db"
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	caseRecord := model.CaseRecord{ID: "case-persist", Title: "Durable matter", ClientID: "client", Category: "class-30", Status: "approved"}
	if err := first.PutCase(caseRecord); err != nil {
		t.Fatal(err)
	}
	if err := first.PutClient(model.ClientProfile{ID: "client", Name: "Ada"}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutReview(model.ReviewRecord{ID: "review-persist", CaseID: caseRecord.ID, Decision: "approved"}); err != nil {
		t.Fatal(err)
	}
	if err := first.PutArchive(model.ArchiveEntry{ID: "archive-persist", CaseID: caseRecord.ID}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	got, err := second.GetCase(caseRecord.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != caseRecord.Title || got.Status != "approved" {
		t.Fatalf("unexpected case: %+v", got)
	}
	reviews, err := second.ListReviews(caseRecord.ID)
	if err != nil || len(reviews) != 1 {
		t.Fatalf("reviews: %v %d", err, len(reviews))
	}
	if _, err := second.GetArchive("archive-persist"); err != nil {
		t.Fatal(err)
	}
}

func TestStoreRejectsInvalidLifecycle(t *testing.T) {
	if _, err := Open(""); err == nil {
		t.Fatal("expected path error")
	}
	s, err := Open(t.TempDir() + "/closed.db")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.PutCase(model.CaseRecord{ID: "x"}); err == nil {
		t.Fatal("expected closed error")
	}
}
