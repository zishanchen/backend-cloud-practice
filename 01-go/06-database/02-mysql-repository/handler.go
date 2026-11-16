package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type TaskHandler struct {
	service TaskService
}

func NewTaskHandler(
	service TaskService,
) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

func (h *TaskHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	tasks, err := h.service.List(
		r.Context(),
	)
	if err != nil {
		log.Printf(
			"list tasks failed: %v",
			err,
		)

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

	task, err := h.service.Get(
		r.Context(),
		id,
	)

	if errors.Is(
		err,
		ErrTaskNotFound,
	) {
		writeError(
			w,
			http.StatusNotFound,
			"task not found",
		)
		return
	}

	if err != nil {
		log.Printf(
			"get task failed: %v",
			err,
		)

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

	task, err := h.service.Create(
		r.Context(),
		input,
	)

	if errors.Is(
		err,
		ErrTitleRequired,
	) {
		writeError(
			w,
			http.StatusBadRequest,
			"title is required",
		)
		return
	}

	if err != nil {
		log.Printf(
			"create task failed: %v",
			err,
		)

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
