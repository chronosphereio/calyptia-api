package types

import "time"

// AggregatorMeasurement stores a list of metrics and totals for an aggregator.
type AggregatorMeasurement struct {
	Metrics map[string][]MetricFields `json:"metrics"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// AggregatorMetrics stores a set of AggregatorMeasurement metrics for an aggregator.
type AggregatorMetrics struct {
	Measurements map[string]AggregatorMeasurement `json:"measurements"`
}

// AddMeasurementMetrics appends a set of metrics/totals to a given measurement field on a AggregatorMetrics struct.
func (a *AggregatorMetrics) AddMeasurementMetrics(measurement string, metrics, totals []MetricFields) {
	if a.Measurements == nil {
		a.Measurements = map[string]AggregatorMeasurement{}
	}
	if _, ok := a.Measurements[measurement]; !ok {
		a.Measurements[measurement] = AggregatorMeasurement{
			Metrics: map[string][]MetricFields{},
			Totals:  map[string][]MetricFields{},
		}
	}
	for _, m := range metrics {
		a.Measurements[measurement].Metrics[m.Field] = append(a.Measurements[measurement].Metrics[m.Field], m)
	}
	for _, t := range totals {
		a.Measurements[measurement].Totals[t.Field] = append(a.Measurements[measurement].Totals[t.Field], t)
	}
}

// ProjectMetrics response payload for project level metrics.
type ProjectMetrics struct {
	Measurements map[string]ProjectMeasurement `json:"measurements"`
	TopPlugins   MetricsPluginTotal            `json:"topPlugins"`
}

// ProjectMeasurement struct to store project metrics, used for project level metrics.
type ProjectMeasurement struct {
	Totals  MeasurementTotal   `json:"totals"`
	Plugins map[string]Metrics `json:"plugins"`
}

// AddMeasurementMetrics appends a set of metrics/totals to a given measurement field on a ProjectMetrics struct.
func (a *ProjectMetrics) AddMeasurementMetrics(measurement string, metrics []MetricFields, total MeasurementTotal) {
	if a.Measurements == nil {
		a.Measurements = map[string]ProjectMeasurement{}
	}

	mm := a.Measurements[measurement]
	if mm.Plugins == nil {
		mm.Plugins = map[string]Metrics{}
	}
	mm.Totals = total

	for _, metric := range metrics {
		plugin := mm.Plugins[metric.Plugin]
		if plugin.Metrics == nil {
			plugin.Metrics = map[string][]MetricFields{}
		}
		plugin.Metrics[metric.Field] = append(plugin.Metrics[metric.Field], metric)
		mm.Plugins[metric.Plugin] = plugin
	}

	a.Measurements[measurement] = mm
}

// AgentMetrics response payload for agent level metrics.
type AgentMetrics struct {
	Measurements map[string]AgentMeasurement `json:"measurements"`
}

// AddMeasurementMetrics appends a set of metrics/totals to a given measurement field on a AgentMetrics struct.
func (a *AgentMetrics) AddMeasurementMetrics(measurement string, metrics []MetricFields, totals []MetricFields) {
	if a.Measurements == nil {
		a.Measurements = map[string]AgentMeasurement{}
	}

	if _, ok := a.Measurements[measurement]; !ok {
		a.Measurements[measurement] = AgentMeasurement{
			Plugins: map[string]Metrics{},
			Totals:  map[string][]MetricFields{},
		}
	}

	for _, m := range metrics {
		plugin := a.Measurements[measurement].Plugins[m.Plugin]
		if plugin.Metrics == nil {
			plugin.Metrics = map[string][]MetricFields{}
		}
		plugin.Metrics[m.Field] = append(plugin.Metrics[m.Field], m)
		a.Measurements[measurement].Plugins[m.Plugin] = plugin
	}

	for _, t := range totals {
		a.Measurements[measurement].Totals[t.Field] = append(a.Measurements[measurement].Totals[t.Field], t)
	}
}

// PipelineMetric response payload for pipeline level metric.
type PipelineMetric struct {
	Data  AgentMetrics `json:"data"`
	Error error        `json:"error"`
}

// AgentMeasurement stores per plugin and total agent level metrics.
type AgentMeasurement struct {
	Plugins map[string]Metrics        `json:"plugins"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// MetricsPluginTotal stores totals per plugin metrics.
type MetricsPluginTotal map[string]map[string]*float64

// MeasurementTotal stores totals per measurement.
type MeasurementTotal map[string]*float64

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

// MetricsParams parameters to filtering metrics by.
type MetricsParams struct {
	Start    time.Duration
	Interval time.Duration
}

// CreatedAgentMetrics response model for created agent metrics.
type CreatedAgentMetrics struct {
	Total uint64 `json:"totalInserted"`
}

// PipelinesMetrics response payload for aggregator level pipeline metrics.
// This type is a map of PipelineID(s) -> PipelineMetric.
type PipelinesMetrics map[string]PipelineMetric

// TODO: define "add metrics" type.
