package workflow30

import (
	"lawindex/internal/archive"
	"lawindex/internal/catalog"
	"lawindex/internal/review"
	"lawindex/internal/store"
	"testing"
)

func TestWorkflow30BusinessInvariant(t *testing.T) {
	storage, err := store.Open(t.TempDir() + "/cases.db")
	if err != nil {
		t.Fatal(err)
	}
	defer storage.Close()
	catalogue := catalog.New(storage)
	reviews := review.New(storage, nil)
	archives := archive.New(storage)
	chain := NewChain(catalogue, reviews, archives)
	caseRecord, err := chain.Intake("Lease dispute", "client-1", "class-30", "Lease renewal evidence")
	if err != nil {
		t.Fatal(err)
	}
	if err := chain.Submit(caseRecord.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := chain.Decide(caseRecord.ID, "attorney", "approved"); err != nil {
		t.Fatal(err)
	}
	entry, err := chain.Close(caseRecord.ID, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if entry.CaseID != caseRecord.ID {
		t.Fatalf("wrong archive case: %s", entry.CaseID)
	}
}
