package workflow30

import (
	"lawindex/internal/api"
	"lawindex/internal/catalog"
	"lawindex/internal/store"
	"net/http/httptest"
	"testing"
)

func TestWorkflowSearchAndReport(t *testing.T) {
	storage, err := store.Open(t.TempDir() + "/search.db")
	if err != nil {
		t.Fatal(err)
	}
	defer storage.Close()
	catalogue := catalog.New(storage)
	if _, err := catalogue.CreateCase("Lease index", "client-a", "class-30", "renewal"); err != nil {
		t.Fatal(err)
	}
	if _, err := catalogue.CreateCase("Trademark index", "client-b", "class-45", "mark"); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest("GET", "/cases?q=lease", nil)
	rec := httptest.NewRecorder()
	api.New(catalogue).Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("search status %d", rec.Code)
	}
	reportReq := httptest.NewRequest("GET", "/report", nil)
	reportRec := httptest.NewRecorder()
	api.New(catalogue).Handler().ServeHTTP(reportRec, reportReq)
	if reportRec.Code != 200 {
		t.Fatalf("report status %d", reportRec.Code)
	}
}
