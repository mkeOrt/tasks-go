package httphandler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/mkeOrt/tasks-go/internal/dto"
)

type mockServiceOptions struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

type mockTaskService struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

func newMockTaskService(so *mockServiceOptions) *mockTaskService {
	return &mockTaskService{getAllFunc: so.getAllFunc}
}

func (m *mockTaskService) GetAll(ctx context.Context) ([]domain.Task, error) {
	return m.getAllFunc(ctx)
}

func TestNewTaskHandler(t *testing.T) {
	h := NewTaskHandler(newMockTaskService(&mockServiceOptions{}))
	if h == nil {
		t.Fatal("expected handler to be initialized")
	}
}

func TestTaskHandler_RegisterRoutes(t *testing.T) {
	h := NewTaskHandler(newMockTaskService(&mockServiceOptions{
		getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
			return []domain.Task{}, nil
		},
	}))
	mux := h.RegisterRoutes()

	req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status OK, got %d", w.Code)
	}
}

func TestTaskHandler_GetAll(t *testing.T) {
	useCases := []struct {
		name           string
		serviceOptions *mockServiceOptions
		expectedStatus int
		expectedBody   []dto.TaskDTO
	}{
		{
			name: "success",
			serviceOptions: &mockServiceOptions{
				getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
					return []domain.Task{
						{ID: 1, Title: "Task 1"},
						{ID: 2, Title: "Task 2"},
					}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody: []dto.TaskDTO{
				{ID: 1, Title: "Task 1", CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"},
				{ID: 2, Title: "Task 2", CreatedAt: "0001-01-01T00:00:00Z", UpdatedAt: "0001-01-01T00:00:00Z"},
			},
		},
		{
			name: "success empty list",
			serviceOptions: &mockServiceOptions{
				getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
					return []domain.Task{}, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   []dto.TaskDTO{},
		},
		{
			name: "error",
			serviceOptions: &mockServiceOptions{
				getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
					return nil, errors.New("service error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   nil,
		},
	}

	for _, uc := range useCases {
		t.Run(uc.name, func(t *testing.T) {
			svc := newMockTaskService(uc.serviceOptions)
			h := NewTaskHandler(svc)

			req := httptest.NewRequest(http.MethodGet, "/api/tasks", nil)
			w := httptest.NewRecorder()

			h.GetAll(w, req)

			if w.Code != uc.expectedStatus {
				t.Errorf("expected status %d, got %d", uc.expectedStatus, w.Code)
			}

			if uc.expectedStatus == http.StatusOK {
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

				if len(uc.expectedBody) == 0 {
					if tasks != nil && len(tasks.Tasks) != 0 {
						t.Fatalf("expected empty list but got %v", tasks)
					}
					return
				}

				if !reflect.DeepEqual(tasks.Tasks, uc.expectedBody) {
					t.Errorf("expected body %v, got %v", uc.expectedBody, tasks.Tasks)
				}
			}
		})
	}
}
