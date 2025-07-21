package receiver

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/devon-caron/metrifuge/logger"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"sync"
	"time"
)

type LogReceiver struct {
	mu      sync.RWMutex
	clients map[string]chan LogEntry
	logs    []LogEntry
}

var (
	receiver *LogReceiver
	once     sync.Once
	log      *logrus.Logger
)

func GetLogReceiver() *LogReceiver {
	once.Do(func() {
		log = logger.Get()
		receiver = &LogReceiver{
			clients: make(map[string]chan LogEntry),
			logs:    make([]LogEntry, 0),
		}
	})
	return receiver
}

// Add a client connection for broadcasting
func (lr *LogReceiver) AddClient(clientID string) chan LogEntry {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	ch := make(chan LogEntry, 100) // Buffered channel
	lr.clients[clientID] = ch
	return ch
}

// Remove client connection
func (lr *LogReceiver) RemoveClient(clientID string) {
	lr.mu.Lock()
	defer lr.mu.Unlock()

	if ch, exists := lr.clients[clientID]; exists {
		close(ch)
		delete(lr.clients, clientID)
	}
}

// HTTP handler for receiving logs via POST
func (lr *LogReceiver) ReceiveLogsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		lr.handleJSONLogs(w, r)
	case "text/plain", "":
		lr.handleTextLogs(w, r)
	case "application/x-ndjson":
		lr.handleNDJSONLogs(w, r)
	default:
		http.Error(w, "Unsupported content type", http.StatusBadRequest)
	}
}

// Handle JSON log entries
func (lr *LogReceiver) handleJSONLogs(w http.ResponseWriter, r *http.Request) {
	var entry LogEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Log received: %s\n", entry.Message)
}

// Handle plain text logs
func (lr *LogReceiver) handleTextLogs(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     "INFO",
		Message:   string(body),
		Source:    r.Header.Get("X-Log-Source"),
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Log received, body: %v\n", entry.Message)
}

// Handle newline-delimited JSON logs (streaming)
func (lr *LogReceiver) handleNDJSONLogs(w http.ResponseWriter, r *http.Request) {
	scanner := bufio.NewScanner(r.Body)
	count := 0

	for scanner.Scan() {
		var entry LogEntry
		if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
			log.Printf("Error parsing log line: %v", err)
			continue
		}

		if entry.Timestamp == "" {
			entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
		}

		count++
	}

	if err := scanner.Err(); err != nil {
		http.Error(w, "Error reading stream", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = fmt.Fprintf(w, "Processed %d log entries\n", count)
}
