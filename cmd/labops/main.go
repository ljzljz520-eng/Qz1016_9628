package main

import (
	"flag"
	"labops/api"
	"labops/service"
	"labops/storage"
	"labops/workflow"
	"log"
	"net/http"
)

func main() {
	path := flag.String("db", "labops.db", "database path")
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()
	s, e := storage.Open(*path)
	if e != nil {
		log.Fatal(e)
	}
	defer s.Close()
	rs := service.NewRecordService(s)
	q := workflow.NewQuery(s)
	log.Printf("listening on %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, api.NewServer(rs, q).Routes()))
}
