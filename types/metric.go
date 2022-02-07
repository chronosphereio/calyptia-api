package types

import "time"

type AggregatorMeasurement struct {
	Metrics map[string][]MetricFields `json:"metrics"`
	Totals  map[string][]MetricFields `json:"totals"`
}

type AggregatorMetrics struct {
	Measurements map[string]AggregatorMeasurement `json:"measurements"`
}

// ProjectMetrics response payload for project level metrics.
type ProjectMetrics struct {
	Measurements map[string]ProjectMeasurement `json:"measurements"`
	TopPlugins   PluginTotal                   `json:"topPlugins"`
}

// ProjectMeasurement struct to store project metrics, used for project level metrics.
type ProjectMeasurement struct {
	Totals  MeasurementTotal   `json:"totals"`
	Plugins map[string]Metrics `json:"plugins"`
}

// AgentMetrics response payload for agent level metrics.
type AgentMetrics struct {
	Measurements map[string]AgentMeasurement `json:"measurements"`
}

// AgentMeasurement stores per plugin and total agent level metrics.
type AgentMeasurement struct {
	Plugins map[string]Metrics        `json:"plugins"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// PluginTotal stores totals per plugin metrics.
type PluginTotal map[string]map[string]*float64

// MeasurementTotal stores totals per measurement.
type MeasurementTotal map[string]*float64

// Metrics stores a dict of metric type and its fields.
type Metrics struct {
	Metrics map[string][]MetricFields `json:"metrics"`
}

// MetricFields stores a tuple of time, value per metric.
type MetricFields struct {
	Time  time.Time `json:"time"`
	Value *float64  `json:"value"`
}

// MetricsParams parameters to filtering metrics by.
type MetricsParams struct {
	Start    time.Duration
	Interval time.Duration
}

// TODO: define "add metrics" type.
