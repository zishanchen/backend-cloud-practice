package main

import (
	"context"
	"errors"
	"testing"
)

type fakeTaskRepository struct {
	createFn func(
		context.Context,
		CreateTaskRequest,
	) (Task, error)

	getFn func(
		context.Context,
		int64,
	) (Task, error)

	listFn func(
		context.Context,
	) ([]Task, error)

	updateDoneFn func(
		context.Context,
		int64,
		bool,
	) (Task, error)
}

func (f *fakeTaskRepository) Create(
	ctx context.Context,
	input CreateTaskRequest,
) (Task, error) {
	if f.createFn != nil {
		return f.createFn(
			ctx,
			input,
		)
	}

	return Task{}, nil
}

func (f *fakeTaskRepository) Get(
	ctx context.Context,
	id int64,
) (Task, error) {
	if f.getFn != nil {
		return f.getFn(
			ctx,
			id,
		)
	}

	return Task{}, nil
}

func (f *fakeTaskRepository) List(
	ctx context.Context,
) ([]Task, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}

	return nil, nil
}

func TestTaskServiceCreate(
	t *testing.T,
) {
	repoErr := errors.New(
		"repository error",
	)

	tests := []struct {
		name           string
		input          CreateTaskRequest
		repoTask       Task
		repoErr        error
		wantTask       Task
		wantErr        error
		wantRepoCalled bool
		wantRepoTitle  string
	}{
		{
			name: "valid task",
			input: CreateTaskRequest{
				Title: "Learn testing",
				Done:  false,
			},
			repoTask: Task{
				ID:    1,
				Title: "Learn testing",
				Done:  false,
			},
			wantTask: Task{
				ID:    1,
				Title: "Learn testing",
				Done:  false,
			},
			wantRepoCalled: true,
			wantRepoTitle:  "Learn testing",
		},
		{
			name: "trim title",
			input: CreateTaskRequest{
				Title: "  Learn Go  ",
			},
			repoTask: Task{
				ID:    2,
				Title: "Learn Go",
			},
			wantTask: Task{
				ID:    2,
				Title: "Learn Go",
			},
			wantRepoCalled: true,
			wantRepoTitle:  "Learn Go",
		},
		{
			name: "empty title",
			input: CreateTaskRequest{
				Title: "   ",
			},
			wantErr:        ErrTitleRequired,
			wantRepoCalled: false,
		},
		{
			name: "repository failure",
			input: CreateTaskRequest{
				Title: "Learn SQL",
			},
			repoErr:        repoErr,
			wantErr:        repoErr,
			wantRepoCalled: true,
			wantRepoTitle:  "Learn SQL",
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				repoCalled := false

				repo := &fakeTaskRepository{
					createFn: func(
						ctx context.Context,
						input CreateTaskRequest,
					) (Task, error) {
						repoCalled = true

						if input.Title !=
							tc.wantRepoTitle {
							t.Fatalf(
								"expected title %q, got %q",
								tc.wantRepoTitle,
								input.Title,
							)
						}

						return tc.repoTask,
							tc.repoErr
					},
				}

				service :=
					NewTaskService(repo)

				got, err :=
					service.Create(
						context.Background(),
						tc.input,
					)

				if tc.wantErr != nil {
					if !errors.Is(
						err,
						tc.wantErr,
					) {
						t.Fatalf(
							"expected error %v, got %v",
							tc.wantErr,
							err,
						)
					}
				} else if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if got != tc.wantTask {
					t.Fatalf(
						"expected task %+v, got %+v",
						tc.wantTask,
						got,
					)
				}

				if repoCalled !=
					tc.wantRepoCalled {
					t.Fatalf(
						"expected repo called=%v, got %v",
						tc.wantRepoCalled,
						repoCalled,
					)
				}
			},
		)
	}
}

func TestTaskServiceGet(
	t *testing.T,
) {
	tests := []struct {
		name     string
		repoTask Task
		repoErr  error
		wantTask Task
		wantErr  error
	}{
		{
			name: "task found",
			repoTask: Task{
				ID:    1,
				Title: "Learn Go",
			},
			wantTask: Task{
				ID:    1,
				Title: "Learn Go",
			},
		},
		{
			name:    "task not found",
			repoErr: ErrTaskNotFound,
			wantErr: ErrTaskNotFound,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				repo := &fakeTaskRepository{
					getFn: func(
						context.Context,
						int64,
					) (Task, error) {
						return tc.repoTask,
							tc.repoErr
					},
				}

				service :=
					NewTaskService(repo)

				got, err :=
					service.Get(
						context.Background(),
						1,
					)

				if tc.wantErr != nil {
					if !errors.Is(
						err,
						tc.wantErr,
					) {
						t.Fatalf(
							"expected %v, got %v",
							tc.wantErr,
							err,
						)
					}
				} else if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if got != tc.wantTask {
					t.Fatalf(
						"expected %+v, got %+v",
						tc.wantTask,
						got,
					)
				}
			},
		)
	}
}

func TestTaskServiceList(
	t *testing.T,
) {
	repoErr := errors.New(
		"repository error",
	)

	tests := []struct {
		name      string
		repoTasks []Task
		repoErr   error
		wantCount int
		wantErr   error
	}{
		{
			name: "tasks returned",
			repoTasks: []Task{
				{ID: 1, Title: "A"},
				{ID: 2, Title: "B"},
			},
			wantCount: 2,
		},
		{
			name:    "repository failure",
			repoErr: repoErr,
			wantErr: repoErr,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				repo := &fakeTaskRepository{
					listFn: func(
						context.Context,
					) ([]Task, error) {
						return tc.repoTasks,
							tc.repoErr
					},
				}

				service :=
					NewTaskService(repo)

				got, err :=
					service.List(
						context.Background(),
					)

				if tc.wantErr != nil {
					if !errors.Is(
						err,
						tc.wantErr,
					) {
						t.Fatalf(
							"expected %v, got %v",
							tc.wantErr,
							err,
						)
					}

					return
				}

				if err != nil {
					t.Fatalf(
						"unexpected error: %v",
						err,
					)
				}

				if len(got) !=
					tc.wantCount {
					t.Fatalf(
						"expected %d tasks, got %d",
						tc.wantCount,
						len(got),
					)
				}
			},
		)
	}
}

func (f *fakeTaskRepository) UpdateDone(
	ctx context.Context,
	id int64,
	done bool,
) (Task, error) {
	if f.updateDoneFn != nil {
		return f.updateDoneFn(
			ctx,
			id,
			done,
		)
	}

	return Task{}, nil
}
