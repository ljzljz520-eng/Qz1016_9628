package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"labops/model"
)

type StatusResponse struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

func writeJSON(w http.ResponseWriter, code int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeRecord(r *http.Request) (model.Record, error) {
	var record model.Record
	if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
		return model.Record{}, err
	}
	record.Title = model.CanonicalTitle(record.Title)
	record.Priority = model.NormalizePriority(record.Priority)
	return record, nil
}

func pathID(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return ""
	}
	return parts[len(parts)-1]
}

func isJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json") || strings.Contains(r.Header.Get("Content-Type"), "application/json")
}

func statusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if strings.Contains(err.Error(), "required") {
		return http.StatusUnprocessableEntity
	}
	if strings.Contains(err.Error(), "denied") {
		return http.StatusConflict
	}
	return http.StatusInternalServerError
}

func methodAllowed(w http.ResponseWriter, allowed ...string) {
	w.Header().Set("Allow", strings.Join(allowed, ", "))
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func requestMethod(r *http.Request) string {
	if r == nil {
		return ""
	}
	return strings.ToUpper(r.Method)
}
