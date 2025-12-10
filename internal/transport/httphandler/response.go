package httphandler

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data"`
	Error   string `json:"error"`
}

func ResponseWithJson(w http.ResponseWriter, statusCode int, data *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func RespondWithJson(w http.ResponseWriter, code int, payload any) {
	ResponseWithJson(w, code, &Response{
		Success: true,
		Data:    payload,
	})
}

func RespondWithErrorJson(w http.ResponseWriter, code int, message string) {
	ResponseWithJson(w, code, &Response{
		Success: false,
		Error:   message,
	})
}
