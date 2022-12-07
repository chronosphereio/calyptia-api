package types

import (
	"time"
)

// AggregatorMeasurement stores a list of metrics and totals for an aggregator.
type AggregatorMeasurement struct {
	Metrics map[string][]MetricFields `json:"metrics"`
	Totals  map[string][]MetricFields `json:"totals"`
}

// CoreInstanceMetricsV1 stores a set of AggregatorMeasurement metrics for an aggregator.
type CoreInstanceMetricsV1 struct {
	Measurements map[string]AggregatorMeasurement `json:"measurements"`
}

// AddMeasurementMetrics appends a set of metrics/totals to a given measurement field on a AggregatorMetrics struct.
func (a *CoreInstanceMetricsV1) AddMeasurementMetrics(measurement string, metrics, totals []MetricFields) {
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

// MetricsSummary stores a list of totals for a core instance.
type MetricsSummary struct {
	Input  MetricsInput  `json:"input"`
	Filter MetricsFilter `json:"filter"`
	Output MetricsOutput `json:"output"`
}

// MetricsInput stores totals for a core instance input.
type MetricsInput struct {
	Bytes   *float64 `json:"bytes"`
	Records *float64 `json:"records"`
}

// MetricsFilter stores totals for a core instance filter.
type MetricsFilter struct {
	DropRecords *float64 `json:"dropRecords"`
	EmitRecords *float64 `json:"emitRecords"`
}

// MetricsOutput stores totals for a core instance output.
type MetricsOutput struct {
	Bytes          *float64 `json:"bytes"`
	Records        *float64 `json:"records"`
	Errors         *float64 `json:"errors"`
	Retries        *float64 `json:"retries"`
	RetriedRecords *float64 `json:"retriedRecords"`
	RetriesFailed  *float64 `json:"retriesFailed"`
	DroppedRecords *float64 `json:"droppedRecords"`
	Loads          *float64 `json:"loads"`
}

// MetricsSummaryPlugin stores a list of totals for a core instance
// for a specific plugin.
type MetricsSummaryPlugin struct {
	Inputs  []MetricsInputPlugin  `json:"inputs"`
	Filters []MetricsFilterPlugin `json:"filters"`
	Outputs []MetricsOutputPlugin `json:"outputs"`
}

type MetricsInputPlugin struct {
	Instance string       `json:"instance"`
	Metrics  MetricsInput `json:"metrics"`
}

type MetricsFilterPlugin struct {
	Instance string        `json:"instance"`
	Metrics  MetricsFilter `json:"metrics"`
}

type MetricsOutputPlugin struct {
	Instance string        `json:"instance"`
	Metrics  MetricsOutput `json:"metrics"`
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
	Error string       `json:"error"`
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

// ToOverTime converts a regular MetricFields to a MetricOverTime.
func (m MetricFields) ToOverTime() MetricOverTime {
	return MetricOverTime{
		Time:  m.Time,
		Value: m.Value,
	}
}

// MetricsParams parameters to filtering metrics by.
type MetricsParams struct {
	Start    time.Duration
	Interval time.Duration
}

// PipelinesMetricsParams request payload for bulk querying pipeline metrics for a given aggregator.
type PipelinesMetricsParams struct {
	MetricsParams
	PipelineIDs []string
}

// CreatedAgentMetrics response model for created agent metrics.
type CreatedAgentMetrics struct {
	Total uint `json:"totalInserted"`
}

// PipelinesMetrics response payload for aggregator level pipeline metrics.
// This type is a map of PipelineID(s) -> PipelineMetric.
type PipelinesMetrics map[string]PipelineMetric

// TODO: define "add metrics" type.

// MetricsOverTime stores a list of metrics over time for a core instance.
type MetricsOverTime struct {
	Input  MetricsOverTimeInput  `json:"input"`
	Filter MetricsOverTimeFilter `json:"filter"`
	Output MetricsOverTimeOutput `json:"output"`
}

// MetricsOverTimeInput stores a list of metrics over time for a core instance input.
type MetricsOverTimeInput struct {
	Bytes   []MetricOverTime `json:"bytes"`
	Records []MetricOverTime `json:"records"`
}

// MetricsOverTimeFilter stores a list of metrics over time for a core instance filter.
type MetricsOverTimeFilter struct {
	DropRecords []MetricOverTime `json:"dropRecords"`
	EmitRecords []MetricOverTime `json:"emitRecords"`
}

// MetricsOverTimeOutput stores a list of metrics over time for a core instance output.
type MetricsOverTimeOutput struct {
	Bytes          []MetricOverTime `json:"bytes"`
	Records        []MetricOverTime `json:"records"`
	Errors         []MetricOverTime `json:"errors"`
	Retries        []MetricOverTime `json:"retries"`
	RetriedRecords []MetricOverTime `json:"retriedRecords"`
	RetriesFailed  []MetricOverTime `json:"retriesFailed"`
	DroppedRecords []MetricOverTime `json:"droppedRecords"`
	Loads          []MetricOverTime `json:"loads"`
}

// MetricsOverTimeByPlugin stores a list of metrics over time for a core instance
// for a specific plugin.
type MetricsOverTimeByPlugin struct {
	Inputs  []MetricsOverTimeByPluginInput  `json:"inputs"`
	Filters []MetricsOverTimeByPluginFilter `json:"filters"`
	Outputs []MetricsOverTimeByPluginOutput `json:"outputs"`
}

// MetricsOverTimeByPluginInput stores a list of metrics over time for core instance inputs
// organized by plugin.
type MetricsOverTimeByPluginInput struct {
	Instance string               `json:"instance"`
	Metrics  MetricsOverTimeInput `json:"metrics"`
}

// MetricsOverTimeByPluginFilter stores a list of metrics over time for core instance filters
// organized by plugin.
type MetricsOverTimeByPluginFilter struct {
	Instance string                `json:"instance"`
	Metrics  MetricsOverTimeFilter `json:"metrics"`
}

// MetricsOverTimeByPluginOutput stores a list of metrics over time for core instance outputs
// organized by plugin.
type MetricsOverTimeByPluginOutput struct {
	Instance string                `json:"instance"`
	Metrics  MetricsOverTimeOutput `json:"metrics"`
}

func (o *MetricsOverTimeOutput) Init() *MetricsOverTimeOutput {
	o.Bytes = []MetricOverTime{}
	o.Records = []MetricOverTime{}
	o.Errors = []MetricOverTime{}
	o.Retries = []MetricOverTime{}
	o.RetriedRecords = []MetricOverTime{}
	o.RetriesFailed = []MetricOverTime{}
	o.DroppedRecords = []MetricOverTime{}
	o.Loads = []MetricOverTime{}
	return o
}

func (f *MetricsOverTimeFilter) Init() *MetricsOverTimeFilter {
	f.DropRecords = []MetricOverTime{}
	f.EmitRecords = []MetricOverTime{}
	return f
}

func (i *MetricsOverTimeInput) Init() *MetricsOverTimeInput {
	i.Bytes = []MetricOverTime{}
	i.Records = []MetricOverTime{}
	return i
}

type MetricOverTime struct {
	Time  time.Time `json:"time"`
	Value *float64  `json:"value"`
}
