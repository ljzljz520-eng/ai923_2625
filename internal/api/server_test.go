package api

import (
	"lawindex/internal/catalog"
	"lawindex/internal/store"
	"net/http/httptest"
	"testing"
)

func TestHTTPCaseSearch(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/api.db")
	defer s.Close()
	c := catalog.New(s)
	_, _ = c.CreateCase("Matter", "c", "class-30", "summary")
	req := httptest.NewRequest("GET", "/cases?q=matter", nil)
	rec := httptest.NewRecorder()
	New(c).Handler().ServeHTTP(rec, req)
	if rec.Code != 200 {
		t.Fatalf("status %d", rec.Code)
	}
}
