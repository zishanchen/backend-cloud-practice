package main

type Task struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type CreateTaskRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
