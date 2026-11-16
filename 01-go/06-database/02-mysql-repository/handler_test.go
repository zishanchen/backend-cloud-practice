package main

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeTaskService struct {
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
}

func (f *fakeTaskService) Create(
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

func (f *fakeTaskService) Get(
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

func (f *fakeTaskService) List(
	ctx context.Context,
) ([]Task, error) {
	if f.listFn != nil {
		return f.listFn(ctx)
	}

	return nil, nil
}

func TestTaskHandlerCreate(
	t *testing.T,
) {
	serviceErr := errors.New(
		"service failure",
	)

	tests := []struct {
		name              string
		body              string
		serviceTask       Task
		serviceErr        error
		wantStatus        int
		wantBodyContains  string
		wantServiceCalled bool
	}{
		{
			name: "created",
			body: `{
				"title":"Learn testing",
				"done":false
			}`,
			serviceTask: Task{
				ID:    1,
				Title: "Learn testing",
				Done:  false,
			},
			wantStatus:        http.StatusCreated,
			wantBodyContains:  `"id":1`,
			wantServiceCalled: true,
		},
		{
			name: "invalid JSON",
			body: `{
				"title":
			}`,
			wantStatus:        http.StatusBadRequest,
			wantBodyContains:  `"error":"invalid JSON"`,
			wantServiceCalled: false,
		},
		{
			name: "empty title",
			body: `{
				"title":""
			}`,
			serviceErr:        ErrTitleRequired,
			wantStatus:        http.StatusBadRequest,
			wantBodyContains:  `"error":"title is required"`,
			wantServiceCalled: true,
		},
		{
			name: "service error",
			body: `{
				"title":"Learn Go"
			}`,
			serviceErr:        serviceErr,
			wantStatus:        http.StatusInternalServerError,
			wantBodyContains:  `"error":"failed to create task"`,
			wantServiceCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				serviceCalled := false

				service := &fakeTaskService{
					createFn: func(
						context.Context,
						CreateTaskRequest,
					) (Task, error) {
						serviceCalled = true

						return tc.serviceTask,
							tc.serviceErr
					},
				}

				handler :=
					NewTaskHandler(service)

				request :=
					httptest.NewRequest(
						http.MethodPost,
						"/tasks",
						bytes.NewBufferString(
							tc.body,
						),
					)

				request.Header.Set(
					"Content-Type",
					"application/json",
				)

				response :=
					httptest.NewRecorder()

				handler.Create(
					response,
					request,
				)

				if response.Code !=
					tc.wantStatus {
					t.Fatalf(
						"expected status %d, got %d",
						tc.wantStatus,
						response.Code,
					)
				}

				if !strings.Contains(
					response.Body.String(),
					tc.wantBodyContains,
				) {
					t.Fatalf(
						"expected body to contain %q, got %q",
						tc.wantBodyContains,
						response.Body.String(),
					)
				}

				if serviceCalled !=
					tc.wantServiceCalled {
					t.Fatalf(
						"expected service called=%v, got %v",
						tc.wantServiceCalled,
						serviceCalled,
					)
				}
			},
		)
	}
}

func TestTaskHandlerGet(
	t *testing.T,
) {
	serviceErr := errors.New(
		"service failure",
	)

	tests := []struct {
		name              string
		id                string
		serviceTask       Task
		serviceErr        error
		wantStatus        int
		wantBodyContains  string
		wantServiceCalled bool
	}{
		{
			name: "task found",
			id:   "1",
			serviceTask: Task{
				ID:    1,
				Title: "Learn Go",
			},
			wantStatus:        http.StatusOK,
			wantBodyContains:  `"id":1`,
			wantServiceCalled: true,
		},
		{
			name:              "invalid id",
			id:                "abc",
			wantStatus:        http.StatusBadRequest,
			wantBodyContains:  `"error":"invalid task id"`,
			wantServiceCalled: false,
		},
		{
			name:              "task not found",
			id:                "999",
			serviceErr:        ErrTaskNotFound,
			wantStatus:        http.StatusNotFound,
			wantBodyContains:  `"error":"task not found"`,
			wantServiceCalled: true,
		},
		{
			name:              "service error",
			id:                "1",
			serviceErr:        serviceErr,
			wantStatus:        http.StatusInternalServerError,
			wantBodyContains:  `"error":"failed to get task"`,
			wantServiceCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				serviceCalled := false

				service := &fakeTaskService{
					getFn: func(
						context.Context,
						int64,
					) (Task, error) {
						serviceCalled = true

						return tc.serviceTask,
							tc.serviceErr
					},
				}

				handler :=
					NewTaskHandler(service)

				request :=
					httptest.NewRequest(
						http.MethodGet,
						"/tasks/"+tc.id,
						nil,
					)

				request.SetPathValue(
					"id",
					tc.id,
				)

				response :=
					httptest.NewRecorder()

				handler.Get(
					response,
					request,
				)

				if response.Code !=
					tc.wantStatus {
					t.Fatalf(
						"expected status %d, got %d",
						tc.wantStatus,
						response.Code,
					)
				}

				if !strings.Contains(
					response.Body.String(),
					tc.wantBodyContains,
				) {
					t.Fatalf(
						"expected body to contain %q, got %q",
						tc.wantBodyContains,
						response.Body.String(),
					)
				}

				if serviceCalled !=
					tc.wantServiceCalled {
					t.Fatalf(
						"expected service called=%v, got %v",
						tc.wantServiceCalled,
						serviceCalled,
					)
				}
			},
		)
	}
}

func TestTaskHandlerList(
	t *testing.T,
) {
	serviceErr := errors.New(
		"service failure",
	)

	tests := []struct {
		name             string
		serviceTasks     []Task
		serviceErr       error
		wantStatus       int
		wantBodyContains string
	}{
		{
			name: "list tasks",
			serviceTasks: []Task{
				{
					ID:    1,
					Title: "Learn testing",
				},
			},
			wantStatus:       http.StatusOK,
			wantBodyContains: `"Learn testing"`,
		},
		{
			name:             "service failure",
			serviceErr:       serviceErr,
			wantStatus:       http.StatusInternalServerError,
			wantBodyContains: `"error":"failed to list tasks"`,
		},
	}

	for _, tc := range tests {
		t.Run(
			tc.name,
			func(t *testing.T) {
				service := &fakeTaskService{
					listFn: func(
						context.Context,
					) ([]Task, error) {
						return tc.serviceTasks,
							tc.serviceErr
					},
				}

				handler :=
					NewTaskHandler(service)

				request :=
					httptest.NewRequest(
						http.MethodGet,
						"/tasks",
						nil,
					)

				response :=
					httptest.NewRecorder()

				handler.List(
					response,
					request,
				)

				if response.Code !=
					tc.wantStatus {
					t.Fatalf(
						"expected status %d, got %d",
						tc.wantStatus,
						response.Code,
					)
				}

				if !strings.Contains(
					response.Body.String(),
					tc.wantBodyContains,
				) {
					t.Fatalf(
						"expected body to contain %q, got %q",
						tc.wantBodyContains,
						response.Body.String(),
					)
				}
			},
		)
	}
}
