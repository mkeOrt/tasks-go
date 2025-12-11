package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/mkeOrt/tasks-go/internal/transport/dto"
)

type mockTaskService struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

func (m *mockTaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	return m.getAllFunc(ctx)
}

func TestNewTaskHandler(t *testing.T) {
	svc := &mockTaskService{
		getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return []domain.Task{}, nil
		},
	}
	h := NewTaskHandler(slog.Default(), svc)
	if h == nil {
		t.Fatal("expected handler to be initialized")
	}
}

func TestTaskHandler_RegisterRoutes(t *testing.T) {
	svc := &mockTaskService{
		getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return []domain.Task{}, nil
		},
	}
	h := NewTaskHandler(slog.Default(), svc)
	mux := h.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}

func TestTaskHandler_GetAll(t *testing.T) {
	testCases := []struct {
		name           string
		setup          func() *mockTaskService
		expectedStatus int
		expectedBody   []dto.TaskDTO
	}{
		{
			name: "should return tasks successfully",
			setup: func() *mockTaskService {
				return &mockTaskService{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return []domain.Task{
							{ID: 1, Title: "Task 1"},
							{ID: 2, Title: "Task 2"},
						}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody: []dto.TaskDTO{
				{ID: 1, Title: "Task 1", CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"},
				{ID: 2, Title: "Task 2", CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"},
			},
		},
		{
			name: "should return empty list successfully",
			setup: func() *mockTaskService {
				return &mockTaskService{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return []domain.Task{}, nil
					},
				}
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []dto.TaskDTO{},
		},
		{
			name: "should return error when service fails",
			setup: func() *mockTaskService {
				return &mockTaskService{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return nil, errors.New("service error")
					},
				}
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			svc := tc.setup()
			h := NewTaskHandler(slog.Default(), svc)

			req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
			w := httptest.NewRecorder()

			h.GetAll(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedStatus == http.StatusOK {
				if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
					t.Errorf("expected Content-Type application/json, got %s", contentType)
				}

				var resp struct {
					Success bool               `json:"success"`
					Data    *dto.TasksResponse `json:"data"`
					Error   string             `json:"error"`
				}
				if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
					t.Fatal(err)
				}

				tasks := resp.Data

				if len(tc.expectedBody) == 0 {
					if tasks != nil && len(tasks.Tasks) != 0 {
						t.Fatalf("expected empty list but got %v", tasks)
					}
					return
				}

				if !reflect.DeepEqual(tasks.Tasks, tc.expectedBody) {
					t.Errorf("expected body %v, got %v", tc.expectedBody, tasks.Tasks)
				}
			}
		})
	}
}
