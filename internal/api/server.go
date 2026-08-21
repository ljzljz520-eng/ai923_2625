package api

import (
	"encoding/json"
	"lawindex/internal/catalog"
	"lawindex/internal/model"
	"lawindex/internal/report"
	"net/http"
)

type Server struct{ catalog *catalog.Service }

func New(c *catalog.Service) *Server { return &Server{catalog: c} }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	s.Register(mux)
	return mux
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusNoContent) }
func (s *Server) cases(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	cases, err := s.catalog.FindCases(model.SearchFilter{Query: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category"), Status: r.URL.Query().Get("status")})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(cases)
}
func (s *Server) report(w http.ResponseWriter, _ *http.Request) {
	cases, err := s.catalog.FindCases(model.SearchFilter{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_ = json.NewEncoder(w).Encode(report.Build(cases))
}
