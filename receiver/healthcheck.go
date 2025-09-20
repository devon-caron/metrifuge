package receiver

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health check endpoint
func (lr *LogReceiver) HealthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	lr.mu.RLock()
	clientCount := len(lr.clients)
	logCount := len(lr.logs)
	lr.mu.RUnlock()

	response := map[string]any{
		"status":     "healthy",
		"clients":    clientCount,
		"total_logs": logCount,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	_ = json.NewEncoder(w).Encode(response)
}
