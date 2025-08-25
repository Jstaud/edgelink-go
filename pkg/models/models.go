package models

import "time"

type Reading struct {
	DeviceID string `json:"device_id"`
	Timestamp time.Time `json:"timestamp"`
	Metrics   map[string]float64 `json:"metrics"`
	Quality   string `json:"quality"`
}

type DeviceSpec struct {
	ID         string        `yaml:"id"` 
	Type       string        `yaml:"type"` // "fake"
	PollEvery  time.Duration `yaml:"poll_every"` // e.g. "1s"
	Jitter     float64       `yaml:"jitter"` // +/- percent (0.05 = 5%)
	MetricNames []string     `yaml:"metric_names"` // e.g. ["temp_c", "pressure_bar"]
}