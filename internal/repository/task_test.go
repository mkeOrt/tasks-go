package repository

import (
	"database/sql"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mkeOrt/tasks-go/internal/domain"
)

func TestNewTaskRepository(t *testing.T) {
	t.Parallel()
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatal("failed to create mock")
	}
	defer db.Close()

	r := NewTaskRepository(db)
	if r == nil {
		t.Fatal("expected repository to not be nil")
	}
}

func TestProductRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("failed to create mock")
	}
	defer db.Close()

	createdAt := time.Now()
	updatedAt := time.Now()

	useCases := []struct {
		name        string
		mockDB      func()
		expected    []domain.Task
		expectedErr error
	}{
		{
			name: "should return error when query fails",
			mockDB: func() {
				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnError(sql.ErrConnDone)
			},
			expected:    nil,
			expectedErr: domain.ErrTaskQueryFailed,
		},
		{
			name: "should return empty list when db returns no rows",
			mockDB: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:    []domain.Task{},
			expectedErr: nil,
		},
		{
			name: "should return populated list when db returns rows",
			mockDB: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"}).
					AddRow(1, "Task 1", false, createdAt, updatedAt).
					AddRow(2, "Task 2", true, createdAt, updatedAt)

				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected: []domain.Task{
				{ID: 1, Title: "Task 1", Done: false, CreatedAt: createdAt, UpdatedAt: updatedAt},
				{ID: 2, Title: "Task 2", Done: true, CreatedAt: createdAt, UpdatedAt: updatedAt},
			},
			expectedErr: nil,
		},
		{
			name: "should return error when rows iteration fails",
			mockDB: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"}).
					AddRow(1, "Task 1", false, createdAt, updatedAt).
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: domain.ErrTaskQueryFailed,
		},
		{
			name: "should return error when scan fails",
			mockDB: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"}).
					AddRow(1, "Task 1", false, createdAt, "invalid-time")

				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: domain.ErrTaskScanFailed,
		},
	}

	for _, uc := range useCases {
		t.Run(uc.name, func(t *testing.T) {
			uc.mockDB()
			repo := NewTaskRepository(db)
			tasks, err := repo.GetAll(t.Context())

			if uc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !errors.Is(err, uc.expectedErr) && !strings.Contains(err.Error(), uc.expectedErr.Error()) {
					t.Fatalf("expected error %v but got %v", uc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
			}

			if len(uc.expected) != 0 && !reflect.DeepEqual(tasks, uc.expected) {
				t.Fatalf("expected tasks %v but got %v", uc.expected, tasks)
			}
		})
	}
}
