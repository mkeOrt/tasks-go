package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

type mockRepositoryOptions struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

type mockTaskRepository struct {
	getAllFunc func(ctx context.Context) ([]domain.Task, error)
}

func newMockTaskRepository(rp *mockRepositoryOptions) *mockTaskRepository {
	return &mockTaskRepository{getAllFunc: rp.getAllFunc}
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

	useCases := []struct {
		name        string
		repoOptions *mockRepositoryOptions
		expected    []domain.Task
		expectedErr error
	}{
		{
			name: "should return error",
			repoOptions: &mockRepositoryOptions{
				getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
					return nil, domain.ErrTasksRetrieveError
				},
			},
			expected:    nil,
			expectedErr: domain.ErrTasksRetrieveError,
		},
		{
			name: "should return empty tasks",
			repoOptions: &mockRepositoryOptions{
				getAllFunc: func(ctx context.Context) ([]domain.Task, error) {
					return []domain.Task{}, nil
				},
			},
			expected:    []domain.Task{},
			expectedErr: nil,
		},
		{
			name: "should return a populated tasks list",
			repoOptions: &mockRepositoryOptions{
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

	for _, uc := range useCases {
		t.Run(uc.name, func(t *testing.T) {
			s := NewTaskService(newMockTaskRepository(uc.repoOptions))
			tasks, err := s.GetAll(t.Context())

			if uc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !errors.Is(err, uc.expectedErr) {
					t.Fatalf("expected error %v but got %v", uc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
			}

			if !reflect.DeepEqual(tasks, uc.expected) {
				t.Fatalf("expected tasks %v but got %v", uc.expected, tasks)
			}
		})
	}
}
