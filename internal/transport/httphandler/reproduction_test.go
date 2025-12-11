package httphandler

import (
	"errors"
	"fmt"
	"testing"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

func TestMapErrorToResponse_Reproduction(t *testing.T) {
	// Simulate the error from Repository
	dbError := errors.New("sql: no such table: tasksx")
	repoError := fmt.Errorf("repository.TaskRepository.GetAll: %w", errors.Join(domain.ErrTaskQueryFailed, dbError))

	// Simulate the error from Service
	serviceError := fmt.Errorf("%w: %w", domain.ErrTasksRetrieveError, repoError)

	// Test mapping
	code, msg := MapErrorToResponse(serviceError)

	expectedMsg := "No se pudo recuperar la lista de tareas"
	if msg != expectedMsg {
		t.Errorf("Expected message %q, got %q. Code: %d", expectedMsg, msg, code)
	}
}
