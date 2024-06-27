package types

import (
	"time"
)

// Label model used internally for metrics.
type Label struct {
	Name  string `json:"name" msgpack:"name"`
	Value string `json:"value" msgpack:"value"`
}

// CreateMetrics types defined by fluent-bit cmetrics v0.5.x.
type CreateMetrics struct {
	// AgentID is set internally by the server and is not part of the request body.
	AgentID string `json:"-" msgpack:"-"`
	// ProjectID is set internally by the server and is not part of the request body.
	ProjectID string `json:"-" msgpack:"-"`
	// FleetID is set internally by the server and is not part of the request body.
	FleetID *string `json:"-" msgpack:"-"`
	// PipelineID is set internally by the server and is not part of the request body.
	PipelineID *string `json:"-" msgpack:"-"`
	// CoreInstanceID is set internally by the server and is not part of the request body.
	CoreInstanceID *string `json:"-" msgpack:"-"`

	Meta    CreateMetricsMeta `json:"meta" msgpack:"meta"`
	Metrics []Metric          `json:"metrics" msgpack:"metrics"`
}

type CreateMetricsMeta struct {
	// Cmetrics   struct{}                    `json:"cmetrics" msgpack:"cmetrics"`
	// External   struct{}                    `json:"external" msgpack:"external"`
	Processing CreateMetricsMetaProcessing `json:"processing" msgpack:"processing"`
}

type CreateMetricsMetaProcessing struct {
	StaticLabels [][2]string `json:"static_abels" msgpack:"static_labels"` // example: [["foo", "bar"], ["my-label", "my-value"]]
}

type Metric struct {
	Meta   MetricMeta    `json:"meta" msgpack:"meta"`
	Values []MetricValue `json:"values" msgpack:"values"`
}

type MetricMeta struct {
	// AggregationType MetricAggregationType `json:"aggregation_type" msgpack:"aggregation_type"`
	// Labels contains only the label keys.
	Labels []string       `json:"labels" msgpack:"labels"` // ex: ["hostname", "name"]
	Opts   MetricMetaOpts `json:"opts" msgpack:"opts"`
	Type   MetricType     `json:"type" msgpack:"type"`
	// Ver    int            `json:"ver" msgpack:"ver"`
	// Buckets []float64 `json:"buckets" msgpack:"buckets"`
}

type MetricAggregationType int

// https://github.com/fluent/fluent-bit/blob/9d9ac68a2b45a4cedeafbf7c5aba513f494eced6/lib/cmetrics/include/cmetrics/cmetrics.h#L32-L34
const (
	MetricAggregationTypeUnspecified MetricAggregationType = iota
	MetricAggregationTypeDelta
	MetricAggregationTypeCumulative
)

func (t MetricAggregationType) String() string {
	switch t {
	case MetricAggregationTypeUnspecified:
		return "unspecified"
	case MetricAggregationTypeDelta:
		return "delta"
	case MetricAggregationTypeCumulative:
		return "cumulative"
	default:
		return "unknown"
	}
}

func (t MetricAggregationType) GoString() string {
	return t.String()
}

type MetricMetaOpts struct {
	Desc      string `json:"desc" msgpack:"desc"`
	Name      string `json:"name" msgpack:"name"` // ex: bytes_total
	Namespace string `json:"ns" msgpack:"ns"`     // ex: fluentbit
	Subsystem string `json:"ss" msgpack:"ss"`     // ex: input
}

type MetricType int

// https://github.com/fluent/fluent-bit/blob/9d9ac68a2b45a4cedeafbf7c5aba513f494eced6/lib/cmetrics/include/cmetrics/cmetrics.h#L26-L30
const (
	MetricTypeCounter MetricType = iota
	MetricTypeGauge
	MetricTypeHistogram
	MetricTypeSummary
	MetricTypeUntyped
)

func (t MetricType) String() string {
	switch t {
	case MetricTypeCounter:
		return "counter"
	case MetricTypeGauge:
		return "gauge"
	case MetricTypeHistogram:
		return "histogram"
	case MetricTypeSummary:
		return "summary"
	case MetricTypeUntyped:
		return "untyped"
	default:
		return "unknown"
	}
}

func (t MetricType) GoString() string {
	return t.String()
}

type MetricValue struct {
	// Hash int64 `json:"hash" msgpack:"hash"`
	// Labels contains the label values from the metric meta label keys.
	// It should match the length of the metric meta label keys.
	Labels []string `json:"labels" msgpack:"labels"` // ex: ["dummy.0"]
	TS     int64    `json:"ts" msgpack:"ts"`         // nanoseconds
	Value  float64  `json:"value" msgpack:"value"`
	// Histogram struct {} `json:"histogram" msgpack:"histogram"`
	// Summary struct {} `json:"summary" msgpack:"summary"`
}

func (v MetricValue) Time() time.Time {
	return time.Unix(0, v.TS)
}

// CreatedMetrics response model for created agent metrics.
type CreatedMetrics struct {
	TotalInserted uint `json:"totalInserted"`
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

// MetricsSummary is the rate of metrics per second.
// For example, the amount of input records per second,
// or the amount of output bytes per second.
// It returns the current (last) value of the metric.
type MetricsSummary struct {
	Input  MetricsInput  `json:"input"`
	Filter MetricsFilter `json:"filter"`
	Output MetricsOutput `json:"output"`
}

// MetricsInput stores totals for a core instance input.
type MetricsInput struct {
	Bytes   float64 `json:"bytes"`
	Records float64 `json:"records"`
}

// MetricsFilter stores totals for a core instance filter.
type MetricsFilter struct {
	Records     float64 `json:"records"`
	Bytes       float64 `json:"bytes"`
	AddRecords  float64 `json:"addRecords"`
	DropRecords float64 `json:"dropRecords"`
	EmitRecords float64 `json:"emitRecords"`
}

// MetricsOutput stores totals for a core instance output.
type MetricsOutput struct {
	Bytes          float64 `json:"bytes"`
	Records        float64 `json:"records"`
	Errors         float64 `json:"errors"`
	Retries        float64 `json:"retries"`
	RetriedRecords float64 `json:"retriedRecords"`
	RetriesFailed  float64 `json:"retriesFailed"`
	DroppedRecords float64 `json:"droppedRecords"`
}

// MetricsOverTime stores a list of metrics over time for a core instance.
type MetricsOverTime struct {
	Input  MetricsOverTimeInput  `json:"input"`
	Filter MetricsOverTimeFilter `json:"filter"`
	Output MetricsOverTimeOutput `json:"output"`
}

func (m *MetricsOverTime) Init() {
	m.Input.Init()
	m.Filter.Init()
	m.Output.Init()
}

// MetricsOverTimeInput stores a list of metrics over time for a core instance input.
type MetricsOverTimeInput struct {
	Bytes   []MetricOverTime `json:"bytes"`
	Records []MetricOverTime `json:"records"`
}

func (m *MetricsOverTimeInput) Init() {
	if m.Bytes == nil {
		m.Bytes = []MetricOverTime{}
	}
	if m.Records == nil {
		m.Records = []MetricOverTime{}
	}
}

// MetricsOverTimeFilter stores a list of metrics over time for a core instance filter.
type MetricsOverTimeFilter struct {
	Bytes       []MetricOverTime `json:"bytes"`
	Records     []MetricOverTime `json:"records"`
	AddRecords  []MetricOverTime `json:"addRecords"`
	DropRecords []MetricOverTime `json:"dropRecords"`
	EmitRecords []MetricOverTime `json:"emitRecords"`
}

func (m *MetricsOverTimeFilter) Init() {
	if m.Bytes == nil {
		m.Bytes = []MetricOverTime{}
	}
	if m.Records == nil {
		m.Records = []MetricOverTime{}
	}
	if m.AddRecords == nil {
		m.AddRecords = []MetricOverTime{}
	}
	if m.DropRecords == nil {
		m.DropRecords = []MetricOverTime{}
	}
	if m.EmitRecords == nil {
		m.EmitRecords = []MetricOverTime{}
	}
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
}

func (m *MetricsOverTimeOutput) Init() {
	if m.Bytes == nil {
		m.Bytes = []MetricOverTime{}
	}
	if m.Records == nil {
		m.Records = []MetricOverTime{}
	}
	if m.Errors == nil {
		m.Errors = []MetricOverTime{}
	}
	if m.Retries == nil {
		m.Retries = []MetricOverTime{}
	}
	if m.RetriedRecords == nil {
		m.RetriedRecords = []MetricOverTime{}
	}
	if m.RetriesFailed == nil {
		m.RetriesFailed = []MetricOverTime{}
	}
	if m.DroppedRecords == nil {
		m.DroppedRecords = []MetricOverTime{}
	}
}

type MetricOverTime struct {
	Time  time.Time `json:"time"`
	Value float64   `json:"value"`
}

// MetricsSummaryPlugin stores a list of totals for a core instance
// for a specific plugin.
type MetricsSummaryPlugin struct {
	Inputs  []MetricsInputPlugin  `json:"inputs"`
	Filters []MetricsFilterPlugin `json:"filters"`
	Outputs []MetricsOutputPlugin `json:"outputs"`
}

func (m *MetricsSummaryPlugin) Init() {
	if m.Inputs == nil {
		m.Inputs = []MetricsInputPlugin{}
	}
	if m.Filters == nil {
		m.Filters = []MetricsFilterPlugin{}
	}
	if m.Outputs == nil {
		m.Outputs = []MetricsOutputPlugin{}
	}
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

// MetricsOverTimeByPlugin stores a list of metrics over time for a core instance
// for a specific plugin.
type MetricsOverTimeByPlugin struct {
	Inputs  []MetricsOverTimeByPluginInput  `json:"inputs"`
	Filters []MetricsOverTimeByPluginFilter `json:"filters"`
	Outputs []MetricsOverTimeByPluginOutput `json:"outputs"`
}

func (m *MetricsOverTimeByPlugin) Init() {
	if m.Inputs == nil {
		m.Inputs = []MetricsOverTimeByPluginInput{}
	}
	if m.Filters == nil {
		m.Filters = []MetricsOverTimeByPluginFilter{}
	}
	if m.Outputs == nil {
		m.Outputs = []MetricsOverTimeByPluginOutput{}
	}
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
