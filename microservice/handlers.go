package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (s *Server) CreateGraph(w http.ResponseWriter, r *http.Request) {
	var in map[string][]float64

	err := json.NewDecoder(r.Body).Decode(&in)
	if err != nil {
		s.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var buf bytes.Buffer

	for k, v := range in {
		if len(v) == 0 {
			http.Error(w, "empty data", http.StatusBadRequest)
			return
		}
		s.Logger.Printf("Creating graph for %s\n", k)
		chart := createLineChart(v)
		err := chart.Render(&buf)
		if err != nil {
			s.Logger.Println("your buff thing broke", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	out := map[string]string{"chart": buf.String()}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(out)
	if err != nil {
		s.Logger.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
