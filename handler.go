package main

import (
	"encoding/json"
	"net/http"
)

func (s *Storage) ListHandler(w http.ResponseWriter, r *http.Request) {
	tasks := s.List()
	json.NewEncoder(w).Encode(tasks)
}

func (s *Storage) CreateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	if input.Text == "" {
		http.Error(w, "поле text обязательно", http.StatusBadRequest)
		return
	}

	task := s.Add(input.Text)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}
