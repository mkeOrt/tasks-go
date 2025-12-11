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

func TestTaskRepository_GetAll(t *testing.T) {
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal("failed to create mock")
	}
	defer db.Close()

	createdAt := time.Now()
	updatedAt := time.Now()

	testCases := []struct {
		name           string
		setup          func()
		expected       []domain.Task
		expectedErr    error
		expectAnyError bool
	}{
		{
			name: "should return error when query fails",
			setup: func() {
				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnError(sql.ErrConnDone)
			},
			expected:    nil,
			expectedErr: sql.ErrConnDone,
		},
		{
			name: "should return empty list when db returns no rows",
			setup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"})
				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:    []domain.Task{},
			expectedErr: nil,
		},
		{
			name: "should return populated list when db returns rows",
			setup: func() {
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
			setup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"}).
					AddRow(1, "Task 1", false, createdAt, updatedAt).
					RowError(0, sql.ErrConnDone)

				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:    nil,
			expectedErr: sql.ErrConnDone,
		},
		{
			name: "should return error when scan fails",
			setup: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "done", "created_at", "updated_at"}).
					AddRow(1, "Task 1", false, createdAt, "invalid-time")

				mock.ExpectQuery("SELECT id, title, done, created_at, updated_at FROM tasks").
					WillReturnRows(rows)
			},
			expected:       nil,
			expectAnyError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.setup()
			repo := NewTaskRepository(db)
			tasks, err := repo.GetAll(t.Context())

			if tc.expectAnyError {
				if err == nil {
					t.Fatal("expected an error but got nil")
				}
			} else if tc.expectedErr != nil {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if !errors.Is(err, tc.expectedErr) && !strings.Contains(err.Error(), tc.expectedErr.Error()) {
					t.Fatalf("expected error %v but got %v", tc.expectedErr, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error but got %v", err)
				}
			}

			if len(tc.expected) != 0 && !reflect.DeepEqual(tasks, tc.expected) {
				t.Fatalf("expected tasks %v but got %v", tc.expected, tasks)
			}
		})
	}
}
