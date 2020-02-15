package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *stats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(s.safe)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func runUIServer(addr string, stats *stats) {
	http.HandleFunc("/stats", stats)
	log.Fatal(http.ListenAndServe(addr, nil))
}
