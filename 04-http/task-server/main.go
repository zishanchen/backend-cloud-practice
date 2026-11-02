package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn net/http", Done: true},
	{ID: 2, Title: "Build /health", Done: true},
	{ID: 3, Title: "Build /tasks", Done: true},
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(tasks)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/tasks", tasksHandler)

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
