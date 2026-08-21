package api

import (
	"encoding/json"
	"lawindex/internal/model"
	"lawindex/internal/report"
	"net/http"
)

func (s *Server) method(w http.ResponseWriter, r *http.Request, expected string) bool {
	if r.Method != expected {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}
func (s *Server) writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func (s *Server) handleCaseQuery(w http.ResponseWriter, r *http.Request) {
	if !s.method(w, r, http.MethodGet) {
		return
	}
	query := model.SearchFilter{Query: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category"), Status: r.URL.Query().Get("status")}
	cases, err := s.catalog.FindCases(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	s.writeJSON(w, cases)
}
func (s *Server) handleReportJSON(w http.ResponseWriter, r *http.Request) {
	if !s.method(w, r, http.MethodGet) {
		return
	}
	cases, err := s.catalog.FindCases(model.SearchFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data, err := report.BuildJSON(cases)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(data)
	_, _ = w.Write([]byte("\n"))
}
func (s *Server) handleCSV(w http.ResponseWriter, r *http.Request) {
	if !s.method(w, r, http.MethodGet) {
		return
	}
	cases, err := s.catalog.FindCases(model.SearchFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/csv")
	_, _ = w.Write([]byte(report.BuildCSV(cases)))
}
