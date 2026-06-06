// Package reporter sends collected metrics to the VisiMon API server via HTTP.
package reporter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/Muhammadkafaby/NetStack/agent/pkg/collector"
)

// Reporter handles sending metrics to the API server.
type Reporter struct {
	serverURL string
	apiKey    string
	client    *http.Client
	hostname  string
}

// New creates a new Reporter.
func New(serverURL, apiKey, hostname string) *Reporter {
	return &Reporter{
		serverURL: serverURL,
		apiKey:    apiKey,
		hostname:  hostname,
		client: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:    5,
				IdleConnTimeout: 60 * time.Second,
			},
		},
	}
}

// SendMetrics sends metrics to the API server.
func (r *Reporter) SendMetrics(metrics *collector.Metrics) error {
	payload, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/metrics", r.serverURL)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if r.apiKey != "" {
		req.Header.Set("X-API-Key", r.apiKey)
	}
	req.Header.Set("User-Agent", "VisiMon-Agent/1.0")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send metrics: %w", err)
	}
	defer resp.Body.Close()

	// Drain body to allow connection reuse
	_, _ = io.Copy(io.Discard, resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

// ReportLoop runs the reporting loop on the given interval.
// It stops when the stop channel is closed.
func (r *Reporter) ReportLoop(interval time.Duration, stop chan struct{}) {
	log.Printf("[visimon] Starting report loop — sending to %s every %s", r.serverURL, interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-stop:
			log.Println("[visimon] Report loop stopped")
			return
		case <-ticker.C:
			metrics, err := collector.CollectAll(r.hostname)
			if err != nil {
				log.Printf("[visimon] Failed to collect metrics: %v", err)
				continue
			}

			if err := r.SendMetrics(metrics); err != nil {
				log.Printf("[visimon] Failed to send metrics: %v", err)
				continue
			}
		}
	}
}
