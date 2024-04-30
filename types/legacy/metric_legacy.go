// Package legacy provides types for legacy/deprecated APIs.
package legacy

import "time"

// ProjectMetrics response payload for project level metrics.
type ProjectMetrics struct {
	Measurements map[string]ProjectMeasurement  `json:"measurements"`
	TopPlugins   map[string]map[string]*float64 `json:"topPlugins"`
}

// AgentMetrics response payload for agent level metrics.
type AgentMetrics struct {
	Measurements map[string]AgentMeasurement `json:"measurements"`
}

// PipelineMetric response payload for pipeline level metric.
type PipelineMetric struct {
	Data  AgentMetrics `json:"data"`
	Error string       `json:"error"`
}

// AgentMeasurement stores per plugin and total agent level metrics.
type AgentMeasurement struct {
	Plugins map[string]Metrics        `json:"plugins"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// ProjectMeasurement struct to store project metrics, used for project level metrics.
type ProjectMeasurement struct {
	Totals  map[string]*float64 `json:"totals"`
	Plugins map[string]Metrics  `json:"plugins"`
}

// AggregatorMeasurement stores a list of metrics and totals for an aggregator.
type AggregatorMeasurement struct {
	Metrics map[string][]MetricFields `json:"metrics"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// CoreInstanceMetrics stores a set of AggregatorMeasurement metrics for an aggregator.
type CoreInstanceMetrics struct {
	Measurements map[string]AggregatorMeasurement `json:"measurements"`
}

// Metrics stores a dict of metric type and its fields.
type Metrics struct {
	Metrics map[string][]MetricFields `json:"metrics"`
}

// MetricFields stores a tuple of time, value per metric.
type MetricFields struct {
	Time   time.Time `json:"time"`
	Value  *float64  `json:"value"`
	Field  string    `json:"-"`
	Plugin string    `json:"-"`
}
