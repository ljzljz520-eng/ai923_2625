package api

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Error string `json:"error"`
}

func WriteError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorBody{Error: err.Error()})
}
func ParseBool(value string) bool { return value == "1" || value == "true" || value == "yes" }
func RequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return "anonymous"
}
func Accepted(w http.ResponseWriter) { w.WriteHeader(http.StatusAccepted) }
