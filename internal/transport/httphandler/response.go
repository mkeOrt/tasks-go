package httphandler

import (
	"encoding/json"
	"net/http"
)

// Response is a generic HTTP response wrapper.
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

// ResponseWithJson writes a JSON response and handles encoding errors.
func ResponseWithJson(w http.ResponseWriter, statusCode int, data *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// RespondWithJson writes a successful JSON response.
func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	ResponseWithJson(w, code, &Response{
		Success: true,
		Data:    payload,
	})
}

// RespondWithErrorJson writes an error JSON response.
func RespondWithErrorJson(w http.ResponseWriter, code int, message string) {
	ResponseWithJson(w, code, &Response{
		Success: false,
		Error:   message,
	})
}
