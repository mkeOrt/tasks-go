package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

type mockTaskRepository struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

func (m *mockTaskRepository) GetAll(ctx context.Context) ([]domain.Task, error) {
	return m.getAllFunc(ctx)
}

func TestNewTaskService(t *testing.T) {
	s := NewTaskService(nil)
	if s == nil {
		t.Fatal("expected service to be initialized")
	}
}

func TestTaskService_GetAll(t *testing.T) {
	testCases := []struct {
		name        string
		setup       func() *mockTaskRepository
		expected    []domain.Task
		expectedErr error
	}{
		{
			name: "should return error when repository fails",
			setup: func() *mockTaskRepository {
				return &mockTaskRepository{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return nil, domain.ErrTaskRetrievalFailed
					},
				}
			},
			expected:    nil,
			expectedErr: domain.ErrTaskRetrievalFailed,
		},
		{
			name: "should return empty tasks list",
			setup: func() *mockTaskRepository {
				return &mockTaskRepository{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return []domain.Task{}, nil
					},
				}
			},
			expected:    []domain.Task{},
			expectedErr: nil,
		},
		{
			name: "should return populated tasks list",
			setup: func() *mockTaskRepository {
				return &mockTaskRepository{
					getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
						return []domain.Task{
							{
								ID:        1,
								Title:     "task 1",
								Done:      false,
								CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
								UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							},
						}, nil
					},
				}
			},
			expected: []domain.Task{
				{
					ID:        1,
					Title:     "task 1",
					Done:      false,
					CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					UpdatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				},
			},
			expectedErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := tc.setup()
			svc := NewTaskService(repo)
			tasks, err := svc.GetAll(t.Context())

			if tc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !errors.Is(err, tc.expectedErr) {
					t.Fatalf("expected error %v but got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
			}

			if !reflect.DeepEqual(tasks, tc.expected) {
				t.Fatalf("expected tasks %v but got %v", tc.expected, tasks)
			}
		})
	}
}
