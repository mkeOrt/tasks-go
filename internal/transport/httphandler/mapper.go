package httphandler

import (
	"errors"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

func MapErrorToStatusCode(err error) int {
	switch {
	case errors.Is(err, domain.ErrTasksRetrieveError):
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
