package domain

import "errors"

// Task errors represent domain-level error conditions.
// These are business errors, not implementation details.
var (
	// ErrTaskNotFound indicates that the requested task does not exist.
	ErrTaskNotFound = errors.New("task not found")
	// ErrTaskRetrievalFailed indicates a failure when fetching tasks.
	ErrTaskRetrievalFailed = errors.New("failed to retrieve tasks")
)
