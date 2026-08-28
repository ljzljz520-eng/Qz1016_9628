package api

import (
	"encoding/json"
	"labops/model"
	"labops/service"
	"labops/workflow"
	"net/http"
)

type Server struct {
	Records *service.RecordService
	Query   *workflow.Query
}

func NewServer(r *service.RecordService, q *workflow.Query) *Server {
	return &Server{Records: r, Query: q}
}
func (s *Server) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
func (s *Server) Create(w http.ResponseWriter, r *http.Request) {
	var rec model.Record
	if json.NewDecoder(r.Body).Decode(&rec) != nil {
		http.Error(w, "bad request", 400)
		return
	}
	if e := s.Records.Register(rec, "api"); e != nil {
		http.Error(w, e.Error(), 422)
		return
	}
	json.NewEncoder(w).Encode(rec)
}
func (s *Server) List(w http.ResponseWriter, r *http.Request) {
	rows, e := s.Query.Search(r.URL.Query().Get("q"), r.URL.Query().Get("status"))
	if e != nil {
		http.Error(w, e.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(rows)
}
func (s *Server) Routes() *http.ServeMux {
	m := http.NewServeMux()
	m.HandleFunc("/health", s.Health)
	m.HandleFunc("/records", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			s.Create(w, r)
		} else {
			s.List(w, r)
		}
	})
	return m
}
