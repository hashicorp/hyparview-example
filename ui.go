package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (s *stats) handle(w http.ResponseWriter, r *http.Request) {
	s.lock.RLock()
	body, _ := json.Marshal(s.safe)
	s.lock.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

type sigmaNode struct {
	ID    string `json:"id"`
	Color int32  `json:"color"`
}

type sigmaEdge struct {
	ID     string `json:"id"`
	Source string `json:"source"`
	Target string `json:"target"`
}

type sigma struct {
	Nodes []*sigmaNode `json:"nodes"`
	Edges []*sigmaEdge `json:"edges"`
}

func (s *stats) handleSigma(w http.ResponseWriter, r *http.Request) {
	data := sigma{
		Nodes: []*sigmaNode{},
		Edges: []*sigmaEdge{},
	}

	s.lock.RLock()
	for id, node := range s.safe {
		data.Nodes = append(data.Nodes, &sigmaNode{
			ID:    id,
			Color: node.App % 8,
		})
	}
	s.lock.RUnlock()

	for _, node := range data.Nodes {
		for _, e := range s.safe[node.ID].Active {
			data.Edges = append(data.Edges, &sigmaEdge{
				ID:     node.ID + e,
				Source: node.ID,
				Target: e,
			})
		}
	}

	body, _ := json.Marshal(&data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func handleStatic(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("static", r.URL.Path)
	info, err := os.Stat(path)
	if err != nil {
		w.WriteHeader(404)
		return
	}

	body := make([]byte, info.Size())
	file, err := os.Open(path)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	count, err := file.Read(body)
	if err != nil || int64(count) != info.Size() {
		w.WriteHeader(500)
		return
	}

	var contentType string
	switch filepath.Ext(path) {
	case "html":
		contentType = "text/html"
	case "js":
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(200)
	w.Write(body)
}

func runUIServer(addr string, stats *stats) {
	log.Printf("debug: starting http %s", addr)
	http.HandleFunc("/stats", stats.handle)
	http.HandleFunc("/sigma", stats.handleSigma)
	http.HandleFunc("/", handleStatic)
	log.Fatal(http.ListenAndServe(addr, nil))
}
