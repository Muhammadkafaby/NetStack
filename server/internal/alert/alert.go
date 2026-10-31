// Package alert defines the types for the VisiMon alert engine.
//
// Phase 1 (MVP) focuses on metrics collection and serving. The full alert
// engine — threshold evaluation, debounce windows, and notification sinks
// (webhook, email, Slack) — will be built on top of these types in a
// follow-up phase.
package alert

// Severity levels for alerts.
const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

// Rule defines a threshold rule to evaluate against incoming metrics.
type Rule struct {
	Name      string  `json:"name"`
	Metric    string  `json:"metric"`   // e.g. "cpu.total_percent"
	Operator  string  `json:"operator"` // ">", "<", ">=", "<=", "=="
	Threshold float64 `json:"threshold"`
	Severity  string  `json:"severity"`
}

// Alert represents a triggered rule for a specific host.
type Alert struct {
	Rule      Rule    `json:"rule"`
	Hostname  string  `json:"hostname"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}
