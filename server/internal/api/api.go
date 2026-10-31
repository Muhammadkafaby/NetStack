// Package api provides REST API handlers for VisiMon.
package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Muhammadkafaby/NetStack/server/internal/db"
	"github.com/Muhammadkafaby/NetStack/server/internal/ws"
)

// Server holds dependencies for API handlers.
type Server struct {
	store *db.Store
	hub   *ws.Hub
}

// New creates a new API server instance.
func New(store *db.Store, hub *ws.Hub) *Server {
	return &Server{
		store: store,
		hub:   hub,
	}
}

// --- Middleware ---

// CORSMiddleware adds CORS headers to all responses.
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-API-Key")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs all requests.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[api] %s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

// --- Helpers ---

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// --- Handlers ---

// IncomingMetrics represents the metrics payload from agents.
type IncomingMetrics struct {
	Timestamp int64       `json:"timestamp"`
	Hostname  string      `json:"hostname"`
	OS        string      `json:"os"`
	Platform  string      `json:"platform"`
	Uptime    uint64      `json:"uptime"`
	CPU       interface{} `json:"cpu"`
	Memory    interface{} `json:"memory"`
	Disks     interface{} `json:"disks"`
	Network   interface{} `json:"network"`
	Load      interface{} `json:"load"`
}

// HandleMetricsPost receives metrics from the agent.
func (s *Server) HandleMetricsPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	var metrics IncomingMetrics
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON: "+err.Error())
		return
	}

	if metrics.Hostname == "" {
		writeError(w, http.StatusBadRequest, "hostname is required")
		return
	}

	// Get API key from header
	apiKey := r.Header.Get("X-API-Key")

	// Upsert host
	hostID, err := s.store.UpsertHost(metrics.Hostname, metrics.OS, metrics.Platform, apiKey)
	if err != nil {
		log.Printf("[api] Failed to upsert host: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to register host")
		return
	}

	// Store the raw metrics as JSON
	rawJSON, err := json.Marshal(metrics)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to serialize metrics")
		return
	}

	if err := s.store.StoreMetric(hostID, metrics.Timestamp, string(rawJSON)); err != nil {
		log.Printf("[api] Failed to store metric: %v", err)
		writeError(w, http.StatusInternalServerError, "Failed to store metric")
		return
	}

	// Broadcast to WebSocket clients
	s.hub.Broadcast(rawJSON)

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// HandleGetLatest returns the latest metrics for a host.
func (s *Server) HandleGetLatest(w http.ResponseWriter, r *http.Request) {
	hostname := extractPathParam(r.URL.Path, "/api/v1/metrics/")
	if hostname == "" || strings.Contains(hostname, "/") {
		writeError(w, http.StatusBadRequest, "Invalid hostname")
		return
	}

	record, err := s.store.GetLatestMetric(hostname)
	if err != nil {
		log.Printf("[api] Failed to get latest metric: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if record == nil {
		writeError(w, http.StatusNotFound, "No metrics found for host")
		return
	}

	// Parse the stored JSON
	var data interface{}
	json.Unmarshal([]byte(record.Data), &data)

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"hostname":  hostname,
		"timestamp": record.Timestamp,
		"metrics":   data,
	})
}

// HandleGetHistory returns time-series metrics for a host.
func (s *Server) HandleGetHistory(w http.ResponseWriter, r *http.Request) {
	// Parse path: /api/v1/metrics/{hostname}/history
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/metrics/")
	parts := strings.SplitN(path, "/", 2)
	if len(parts) != 2 || parts[1] != "history" {
		writeError(w, http.StatusBadRequest, "Invalid path")
		return
	}
	hostname := parts[0]

	// Parse range query param
	rangeStr := r.URL.Query().Get("range")
	if rangeStr == "" {
		rangeStr = "5m"
	}

	duration, err := time.ParseDuration(rangeStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid range format (e.g., 5m, 1h, 7d)")
		return
	}

	since := time.Now().Unix() - int64(duration.Seconds())

	records, err := s.store.GetHistory(hostname, since)
	if err != nil {
		log.Printf("[api] Failed to get history: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	// Parse metrics data
	type MetricPoint struct {
		Timestamp int64       `json:"timestamp"`
		Metrics   interface{} `json:"metrics"`
	}

	points := make([]MetricPoint, len(records))
	for i, rec := range records {
		var data interface{}
		json.Unmarshal([]byte(rec.Data), &data)
		points[i] = MetricPoint{
			Timestamp: rec.Timestamp,
			Metrics:   data,
		}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"hostname": hostname,
		"range":    rangeStr,
		"points":   points,
	})
}

// HandleListHosts returns all registered hosts.
func (s *Server) HandleListHosts(w http.ResponseWriter, r *http.Request) {
	hosts, err := s.store.ListHosts()
	if err != nil {
		log.Printf("[api] Failed to list hosts: %v", err)
		writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if hosts == nil {
		hosts = []db.Host{}
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"hosts": hosts,
		"count": len(hosts),
	})
}

// HandleHealth returns health status.
func (s *Server) HandleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"version":   "0.1.0",
		"timestamp": time.Now().Unix(),
	})
}

// extractPathParam extracts the hostname from paths like /api/v1/metrics/{hostname}
func extractPathParam(path, prefix string) string {
	param := strings.TrimPrefix(path, prefix)
	// Remove trailing slash
	param = strings.TrimSuffix(param, "/")
	// Remove any sub-paths like /history
	if idx := strings.Index(param, "/"); idx > 0 {
		param = param[:idx]
	}
	return param
}

// cleanupTask runs periodically to delete old metrics.
func (s *Server) CleanupTask(interval time.Duration, retention time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		deleted, err := s.store.DeleteOldMetrics(retention)
		if err != nil {
			log.Printf("[cleanup] Error: %v", err)
		} else if deleted > 0 {
			log.Printf("[cleanup] Deleted %d old metric records", deleted)
		}
	}
}

// RegisterRoutes sets up all HTTP routes on the given mux.
func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/metrics", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			s.HandleMetricsPost(w, r)
		default:
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
		}
	})

	// GET /api/v1/metrics/{hostname} - latest
	// GET /api/v1/metrics/{hostname}/history - history
	mux.HandleFunc("/api/v1/metrics/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "Method not allowed")
			return
		}

		path := strings.TrimPrefix(r.URL.Path, "/api/v1/metrics/")
		if strings.HasSuffix(path, "/history") {
			s.HandleGetHistory(w, r)
		} else {
			s.HandleGetLatest(w, r)
		}
	})

	mux.HandleFunc("/api/v1/hosts", s.HandleListHosts)
	mux.HandleFunc("/api/v1/health", s.HandleHealth)
	mux.HandleFunc("/api/v1/ws", s.hub.HandleWebSocket)
}

// ParsePort parses the port flag value.
func ParsePort(portStr string) string {
	p, err := strconv.Atoi(portStr)
	if err != nil || p < 1 || p > 65535 {
		return "8080"
	}
	return portStr
}
