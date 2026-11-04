package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn net/http", Done: true},
	{ID: 2, Title: "Build /health", Done: true},
	{ID: 3, Title: "Learn Go 1.22 ServeMux", Done: false},
}

type CreateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

var nextID = 4

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
}

func getTasks(w http.ResponseWriter) {
	writeJSON(
		w,
		http.StatusOK,
		tasks,
	)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	task := Task{
		ID:    nextID,
		Title: input.Title,
		Done:  input.Done,
	}

	nextID++

	tasks = append(tasks, task)

	writeJSON(
		w,
		http.StatusOK,
		task,
	)

}

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(
		w,
		http.StatusOK,
		tasks,
	)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	idText := r.PathValue("id")

	id, err := strconv.Atoi(idText)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			writeJSON(w, http.StatusOK, task)
			return
		}
	}

	writeError(
		w,
		http.StatusNotFound,
		"task not found",
	)
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}
	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		writeError(w, http.StatusBadRequest, "Title is required")
		return
	}

	task := Task{
		ID:    nextID,
		Title: input.Title,
		Done:  input.Done,
	}

	nextID++

	writeJSON(
		w,
		http.StatusOK,
		task,
	)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /tasks", getTasksHandler)
	mux.HandleFunc("GET /tasks/{id}", getTaskHandler)
	mux.HandleFunc("POST /tasks", createTaskHandler)

	fmt.Println("Server running on http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
