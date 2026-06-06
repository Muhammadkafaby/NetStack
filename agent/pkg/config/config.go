// Package config provides YAML-based configuration for the VisiMon agent.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config represents the agent configuration.
type Config struct {
	ServerURL string `yaml:"server_url" json:"server_url"`
	Interval  int    `yaml:"interval" json:"interval"` // seconds
	Hostname  string `yaml:"hostname" json:"hostname"`
	APIKey    string `yaml:"api_key" json:"api_key"`
	LogLevel  string `yaml:"log_level" json:"log_level"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	hostname, _ := os.Hostname()
	return &Config{
		ServerURL: "http://localhost:8080",
		Interval:  1,
		Hostname:  hostname,
		APIKey:    "",
		LogLevel:  "info",
	}
}

// LoadConfig reads a YAML config file and returns a Config.
// If the file doesn't exist, it returns the default config.
func LoadConfig(path string) (*Config, error) {
	cfg := DefaultConfig()

	if path == "" {
		// Look for config in default locations
		locations := []string{
			"/etc/visimon/agent.yaml",
			"/etc/visimon/agent.yml",
			"./visimon.yaml",
			"./visimon.yml",
			"./agent.yaml",
			"./agent.yml",
		}

		for _, loc := range locations {
			if _, err := os.Stat(loc); err == nil {
				path = loc
				break
			}
		}

		if path == "" {
			return cfg, nil // No config file found, use defaults
		}
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	// If no hostname in config, use OS hostname
	if cfg.Hostname == "" {
		hostname, err := os.Hostname()
		if err == nil {
			cfg.Hostname = hostname
		}
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg, nil
}

// ConfigPath returns the resolved config file path for display.
func ConfigPath(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}
