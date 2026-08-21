package api

import (
	"lawindex/internal/model"
	"net/http"
)

func (s *Server) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/cases", s.handleCaseQuery)
	mux.HandleFunc("/report", s.handleReportJSON)
	mux.HandleFunc("/report.csv", s.handleCSV)
}
func (s *Server) Route(path string) string {
	switch path {
	case "/health":
		return "health"
	case "/cases":
		return "cases"
	case "/report":
		return "report"
	case "/report.csv":
		return "csv"
	default:
		return "not-found"
	}
}
func (s *Server) HandlerWithRoutes() http.Handler {
	mux := http.NewServeMux()
	s.Register(mux)
	return ChainMiddleware(mux, nil)
}
func (s *Server) CaseFilter(r *http.Request) model.SearchFilter {
	return model.SearchFilter{Query: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category"), Status: r.URL.Query().Get("status")}
}
func (s *Server) IsKnownRoute(path string) bool { return s.Route(path) != "not-found" }
