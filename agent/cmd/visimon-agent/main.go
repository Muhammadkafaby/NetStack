// VisiMon Agent — main entry point.
//
// Lightweight system metrics collector that sends data to the VisiMon API server.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Muhammadkafaby/NetStack/agent/pkg/config"
	"github.com/Muhammadkafaby/NetStack/agent/pkg/reporter"
)

var version = "0.1.0"

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "Path to config file")
	serverURL := flag.String("server", "", "VisiMon API server URL (overrides config)")
	interval := flag.Int("interval", 0, "Collection interval in seconds (overrides config)")
	agentHostname := flag.String("hostname", "", "Agent hostname (overrides config)")
	apiKey := flag.String("api-key", "", "API key for server auth (overrides config)")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("VisiMon Agent v%s\n", version)
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("[visimon] Failed to load config: %v", err)
	}

	// Command-line overrides
	if *serverURL != "" {
		cfg.ServerURL = *serverURL
	}
	if *interval > 0 {
		cfg.Interval = *interval
	}
	if *agentHostname != "" {
		cfg.Hostname = *agentHostname
	}
	if *apiKey != "" {
		cfg.APIKey = *apiKey
	}

	// Validate config
	if cfg.ServerURL == "" {
		log.Fatal("[visimon] server_url is required")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 1
	}
	if cfg.Hostname == "" {
		hostname, err := os.Hostname()
		if err != nil {
			log.Fatalf("[visimon] Failed to determine hostname: %v", err)
		}
		cfg.Hostname = hostname
	}

	log.Printf("[visimon] VisiMon Agent v%s starting", version)
	log.Printf("[visimon] Server: %s", cfg.ServerURL)
	log.Printf("[visimon] Hostname: %s", cfg.Hostname)
	log.Printf("[visimon] Interval: %ds", cfg.Interval)

	// Create reporter
	r := reporter.New(cfg.ServerURL, cfg.APIKey, cfg.Hostname)

	// Set up signal handling for graceful shutdown
	stop := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		sig := <-sigCh
		log.Printf("[visimon] Received signal %v, shutting down...", sig)
		close(stop)
	}()

	// Run the report loop (blocks until stop)
	r.ReportLoop(time.Duration(cfg.Interval)*time.Second, stop)

	log.Println("[visimon] Agent stopped")
}
