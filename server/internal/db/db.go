// Package db provides SQLite storage for VisiMon metrics.
package db

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "modernc.org/sqlite"
)

// Store handles all database operations.
type Store struct {
	db *sql.DB
}

// Host represents a registered monitoring agent.
type Host struct {
	ID        int64  `json:"id"`
	Hostname  string `json:"hostname"`
	OS        string `json:"os"`
	Platform  string `json:"platform"`
	FirstSeen string `json:"first_seen"`
	LastSeen  string `json:"last_seen"`
	APIKey    string `json:"-"`
}

// MetricRecord represents a stored metrics entry.
type MetricRecord struct {
	ID        int64  `json:"id"`
	HostID    int64  `json:"host_id"`
	Timestamp int64  `json:"timestamp"`
	Data      string `json:"data"`
}

// New opens or creates the SQLite database and initializes tables.
func New(dbPath string) (*Store, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(1) // SQLite doesn't support concurrent writes

	store := &Store{db: db}
	if err := store.init(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	return store, nil
}

func (s *Store) init() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS hosts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			hostname TEXT UNIQUE NOT NULL,
			os TEXT DEFAULT '',
			platform TEXT DEFAULT '',
			first_seen INTEGER NOT NULL,
			last_seen INTEGER NOT NULL,
			api_key TEXT DEFAULT ''
		)`,
		`CREATE TABLE IF NOT EXISTS metrics (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			host_id INTEGER NOT NULL,
			timestamp INTEGER NOT NULL,
			data TEXT NOT NULL,
			FOREIGN KEY (host_id) REFERENCES hosts(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_host_time ON metrics(host_id, timestamp)`,
		`CREATE INDEX IF NOT EXISTS idx_metrics_timestamp ON metrics(timestamp)`,
	}

	for _, q := range queries {
		if _, err := s.db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

// UpsertHost creates or updates a host record.
func (s *Store) UpsertHost(hostname, os, platform, apiKey string) (int64, error) {
	now := time.Now().Unix()

	result, err := s.db.Exec(
		`INSERT INTO hosts (hostname, os, platform, first_seen, last_seen, api_key)
		 VALUES (?, ?, ?, ?, ?, ?)
		 ON CONFLICT(hostname) DO UPDATE SET
		   os=excluded.os,
		   platform=excluded.platform,
		   last_seen=excluded.last_seen,
		   api_key=CASE WHEN excluded.api_key != '' THEN excluded.api_key ELSE hosts.api_key END`,
		hostname, os, platform, now, now, apiKey,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to upsert host: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get host id: %w", err)
	}

	return id, nil
}

// StoreMetric inserts a new metrics record for a host.
func (s *Store) StoreMetric(hostID int64, timestamp int64, dataJSON string) error {
	_, err := s.db.Exec(
		`INSERT INTO metrics (host_id, timestamp, data) VALUES (?, ?, ?)`,
		hostID, timestamp, dataJSON,
	)
	if err != nil {
		return fmt.Errorf("failed to store metric: %w", err)
	}
	return nil
}

// GetLatestMetric returns the most recent metric for a hostname.
func (s *Store) GetLatestMetric(hostname string) (*MetricRecord, error) {
	row := s.db.QueryRow(
		`SELECT m.id, m.host_id, m.timestamp, m.data
		 FROM metrics m
		 JOIN hosts h ON h.id = m.host_id
		 WHERE h.hostname = ?
		 ORDER BY m.timestamp DESC
		 LIMIT 1`,
		hostname,
	)

	record := &MetricRecord{}
	err := row.Scan(&record.ID, &record.HostID, &record.Timestamp, &record.Data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest metric: %w", err)
	}

	return record, nil
}

// GetHistory returns metrics within a time range for a host.
func (s *Store) GetHistory(hostname string, since int64) ([]MetricRecord, error) {
	rows, err := s.db.Query(
		`SELECT m.id, m.host_id, m.timestamp, m.data
		 FROM metrics m
		 JOIN hosts h ON h.id = m.host_id
		 WHERE h.hostname = ? AND m.timestamp >= ?
		 ORDER BY m.timestamp ASC`,
		hostname, since,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get history: %w", err)
	}
	defer rows.Close()

	var records []MetricRecord
	for rows.Next() {
		var r MetricRecord
		if err := rows.Scan(&r.ID, &r.HostID, &r.Timestamp, &r.Data); err != nil {
			return nil, fmt.Errorf("failed to scan metric: %w", err)
		}
		records = append(records, r)
	}

	return records, rows.Err()
}

// ListHosts returns all registered hosts.
func (s *Store) ListHosts() ([]Host, error) {
	rows, err := s.db.Query(
		`SELECT id, hostname, os, platform, first_seen, last_seen
		 FROM hosts ORDER BY last_seen DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list hosts: %w", err)
	}
	defer rows.Close()

	var hosts []Host
	for rows.Next() {
		var h Host
		if err := rows.Scan(&h.ID, &h.Hostname, &h.OS, &h.Platform, &h.FirstSeen, &h.LastSeen); err != nil {
			return nil, fmt.Errorf("failed to scan host: %w", err)
		}
		hosts = append(hosts, h)
	}

	return hosts, rows.Err()
}

// FindHostByAPIKey looks up a host by API key.
func (s *Store) FindHostByAPIKey(apiKey string) (*Host, error) {
	row := s.db.QueryRow(
		`SELECT id, hostname, os, platform, first_seen, last_seen, api_key
		 FROM hosts WHERE api_key = ? LIMIT 1`,
		apiKey,
	)

	host := &Host{}
	err := row.Scan(&host.ID, &host.Hostname, &host.OS, &host.Platform, &host.FirstSeen, &host.LastSeen, &host.APIKey)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to find host by api key: %w", err)
	}

	return host, nil
}

// DeleteOldMetrics removes metrics older than the given retention period.
func (s *Store) DeleteOldMetrics(retention time.Duration) (int64, error) {
	cutoff := time.Now().Unix() - int64(retention.Seconds())

	result, err := s.db.Exec(`DELETE FROM metrics WHERE timestamp < ?`, cutoff)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old metrics: %w", err)
	}

	deleted, _ := result.RowsAffected()
	if deleted > 0 {
		log.Printf("[db] Cleaned up %d old metric records", deleted)
	}

	// Also clean up orphaned hosts with no metrics
	_, _ = s.db.Exec(
		`DELETE FROM hosts WHERE id NOT IN (SELECT DISTINCT host_id FROM metrics)`,
	)

	return deleted, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
