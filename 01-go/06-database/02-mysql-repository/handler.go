package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type TaskHandler struct {
	repo TaskRepository
}

func NewTaskHandler(
	repo TaskRepository,
) *TaskHandler {
	return &TaskHandler{
		repo: repo,
	}
}

func (h *TaskHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	tasks, err := h.repo.List(
		r.Context(),
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to list tasks",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		tasks,
	)
}

func (h *TaskHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	idText := r.PathValue("id")

	id, err := strconv.ParseInt(
		idText,
		10,
		64,
	)
	if err != nil || id <= 0 {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid task id",
		)
		return
	}

	task, err := h.repo.Get(
		r.Context(),
		id,
	)

	if errors.Is(err, ErrTaskNotFound) {
		writeError(
			w,
			http.StatusNotFound,
			"task not found",
		)
		return
	}

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to get task",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		task,
	)
}

func (h *TaskHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var input CreateTaskRequest

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&input); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"invalid JSON",
		)
		return
	}

	input.Title = strings.TrimSpace(
		input.Title,
	)

	if input.Title == "" {
		writeError(
			w,
			http.StatusBadRequest,
			"title is required",
		)
		return
	}

	task, err := h.repo.Create(
		r.Context(),
		input,
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"failed to create task",
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		task,
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(
	w http.ResponseWriter,
	status int,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Error: message,
		},
	)
}
