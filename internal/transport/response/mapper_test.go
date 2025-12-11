package response

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

func TestMapErrorToResponse(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "Tasks Retrieve Error",
			err:            domain.ErrTasksRetrieveError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "No se pudo recuperar la lista de tareas",
		},
		{
			name:           "Unknown Error",
			err:            errors.New("unknown error"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Ocurrió un error inesperado al procesar la solicitud",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, msg := MapErrorToResponse(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, status)
			}
			if msg != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, msg)
			}
		})
	}
}
