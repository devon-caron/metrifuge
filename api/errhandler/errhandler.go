package errhandler

import (
	"encoding/json"
	"net/http"
)

// Error Response
type ErrorResponse struct {
	// HTTP Status Code
	Code int

	// Error msg
	Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
	response := ErrorResponse{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return
	}
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "an unexpected error occurred", http.StatusInternalServerError)
	}
)
