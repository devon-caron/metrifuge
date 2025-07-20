package api

import (
	"encoding/json"
	"net/http"
)

// Coin Balance Params
type CoinBalanceParams struct {
	Username string
}

// Coin Balance Response
type CoinBalanceResponse struct {
	// HTTP Status Code
	Code int

	// Account Balance
	Balance int64
}

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
