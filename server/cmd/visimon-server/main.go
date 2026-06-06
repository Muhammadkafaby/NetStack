// VisiMon API Server — main entry point.
//
// Receives metrics from agents, stores them in SQLite, and serves them
// via REST API and WebSocket for the dashboard.
package main

import (
    "flag"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "strconv"
    "syscall"
    "time"

    "github.com/Muhammadkafaby/NetStack/server/internal/api"
    "github.com/Muhammadkafaby/NetStack/server/internal/db"
    "github.com/Muhammadkafaby/NetStack/server/internal/ws"
)

var version = "0.1.0"

func main() {
    port := flag.String("port", "8080", "HTTP server port")
    dbPath := flag.String("db-path", "./visimon.db", "Path to SQLite database")
    retention := flag.String("retention", "7d", "Data retention period (e.g., 7d, 30d, 90d)")
    cleanupInterval := flag.String("cleanup-interval", "1h", "How often to run cleanup")
    showVersion := flag.Bool("version", false, "Show version and exit")
    flag.Parse()

    if *showVersion {
        fmt.Printf("VisiMon Server v%s\n", version)
        os.Exit(0)
    }

    // Parse retention duration (support "d" suffix for days)
    retentionDur, err := parseDuration(*retention)
    if err != nil {
        log.Fatalf("[server] Invalid retention format '%s': %v", *retention, err)
    }

    cleanupDur, err := parseDuration(*cleanupInterval)
    if err != nil {
        log.Fatalf("[server] Invalid cleanup interval format '%s': %v", *cleanupInterval, err)
    }

    log.Printf("[server] VisiMon Server v%s starting", version)
    log.Printf("[server] Database: %s", *dbPath)
    log.Printf("[server] Retention: %s", retentionDur)
    log.Printf("[server] Port: %s", *port)

    // Initialize database
    store, err := db.New(*dbPath)
    if err != nil {
        log.Fatalf("[server] Failed to initialize database: %v", err)
    }
    defer store.Close()

    // Initialize WebSocket hub
    hub := ws.NewHub()
    go hub.Run()

    // Initialize API handlers
    apiServer := api.New(store, hub)

    // Set up HTTP routes
    mux := http.NewServeMux()
    apiServer.RegisterRoutes(mux)

    // Apply middleware
    handler := api.LoggingMiddleware(api.CORSMiddleware(mux))

    // Start cleanup goroutine
    go apiServer.CleanupTask(cleanupDur, retentionDur)

    // Create server
    server := &http.Server{
        Addr:         "0.0.0.0:" + *port,
        Handler:      handler,
        ReadTimeout:  30 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Graceful shutdown
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        log.Printf("[server] Listening on 0.0.0.0:%s", *port)
        if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("[server] Server error: %v", err)
        }
    }()

    <-stop
    log.Println("[server] Shutting down...")
    server.Close()
    log.Println("[server] Stopped")
}

// parseDuration parses duration strings, supporting "d" suffix for days.
func parseDuration(s string) (time.Duration, error) {
    if len(s) > 1 && s[len(s)-1] == 'd' {
        days, err := strconv.Atoi(s[:len(s)-1])
        if err != nil {
            return 0, fmt.Errorf("invalid days: %s", s)
        }
        return time.Duration(days) * 24 * time.Hour, nil
    }
    return time.ParseDuration(s)
}
