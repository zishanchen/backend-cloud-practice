package main

import (
	"context"
	"errors"
	"strings"
)

var ErrTitleRequired = errors.New("title is required")

type TaskService interface {
	Create(
		ctx context.Context,
		input CreateTaskRequest,
	) (Task, error)

	Get(
		ctx context.Context,
		id int64,
	) (Task, error)

	List(
		ctx context.Context,
	) ([]Task, error)

	UpdateDone(
		ctx context.Context,
		id int64,
		done bool,
	) (Task, error)
}

type taskService struct {
	repo TaskRepository
}

func NewTaskService(
	repo TaskRepository,
) TaskService {
	return &taskService{
		repo: repo,
	}
}

func (s *taskService) Create(
	ctx context.Context,
	input CreateTaskRequest,
) (Task, error) {
	input.Title = strings.TrimSpace(input.Title)

	if input.Title == "" {
		return Task{}, ErrTitleRequired
	}

	return s.repo.Create(
		ctx,
		input,
	)
}

func (s *taskService) Get(
	ctx context.Context,
	id int64,
) (Task, error) {
	return s.repo.Get(
		ctx,
		id,
	)
}

func (s *taskService) List(
	ctx context.Context,
) ([]Task, error) {
	return s.repo.List(ctx)
}

func (s *taskService) UpdateDone(
	ctx context.Context,
	id int64,
	done bool,
) (Task, error) {
	return s.repo.UpdateDone(
		ctx,
		id,
		done,
	)
}
