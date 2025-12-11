package response

import (
	"errors"
	"net/http"

	"github.com/mkeOrt/tasks-go/internal/domain"
)

// MapErrorToResponse maps an error to a status code and a user-friendly message.
// It uses an allowlist approach: only known errors get specific messages.
// Everything else returns a generic internal error message.
func MapErrorToResponse(err error) (int, string) {
	switch {
	case errors.Is(err, domain.ErrTaskRetrievalFailed):
		return http.StatusInternalServerError, ErrMsgTaskRetrieve
	default:
		return http.StatusInternalServerError, ErrMsgUnexpected
	}
}
