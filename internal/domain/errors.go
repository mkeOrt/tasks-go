package domain

import "errors"

// Tasks errors
var (
	ErrTasksRetrieveError = errors.New("failed to retrieve task list")
	ErrTaskQueryFailed    = errors.New("failed to execute query")
	ErrTaskScanFailed     = errors.New("failed to scan task")
)
