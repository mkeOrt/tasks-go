package httphandler

import (
	"errors"
	"net/http"
	"testing"

	"github.com/mkeOrt/tasks-go/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapErrorToStatusCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "Tasks Retrieve Error",
			err:      domain.ErrTasksRetrieveError,
			expected: http.StatusInternalServerError,
		},
		{
			name:     "Unknown Error",
			err:      errors.New("unknown error"),
			expected: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapErrorToStatusCode(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
