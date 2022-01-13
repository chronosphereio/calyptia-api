package types

import "time"

type ProjectMetrics struct {
	Measurements map[string]ProjectMeasurement `json:"measurements"`
	TopPlugins   PluginTotal                   `json:"topPlugins"`
}

type ProjectMeasurement struct {
	Totals  MeasurementTotal   `json:"totals"`
	Plugins map[string]Metrics `json:"plugins"`
}

type AgentMetrics struct {
	Measurements map[string]AgentMeasurement `json:"measurements"`
}

type AgentMeasurement struct {
	Plugins map[string]Metrics        `json:"plugins"`
	Totals  map[string][]MetricFields `json:"totals"`
}

type PluginTotal map[string]map[string]*float64

type MeasurementTotal map[string]*float64

type Metrics struct {
	Metrics map[string][]MetricFields `json:"metrics"`
}

type MetricFields struct {
	Time  time.Time `json:"time"`
	Value *float64  `json:"value"`
}

type MetricsParams struct {
	Start    time.Duration
	Interval time.Duration
}

// TODO: define "add metrics" type.
