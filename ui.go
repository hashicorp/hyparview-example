package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *stats) handle(w http.ResponseWriter, r *http.Request) {
	body, _ := json.Marshal(s.safe)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func runUIServer(addr string, stats *stats) {
	log.Printf("debug: starting http %s", addr)
	http.HandleFunc("/stats", stats.handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}
