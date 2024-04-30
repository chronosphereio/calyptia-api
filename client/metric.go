package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/calyptia/api/types"
)

// CreateMetrics ingests agent metrics via the metrics API.
func (c *Client) CreateMetrics(ctx context.Context, in types.CreateMetrics) (types.CreatedMetrics, error) {
	var out types.CreatedMetrics
	path := "/v1/agent_metrics"

	data, err := msgpack.Marshal(in)
	if err != nil {
		return out, err
	}

	return out, c.do(ctx, http.MethodPost, path, data, &out)
}

// AgentMetrics contains an overview of the Agent metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) AgentMetrics(ctx context.Context, agentID string) (types.MetricsSummary, error) {
	var out types.MetricsSummary
	path := "/v2/agents/" + url.PathEscape(agentID) + "/metrics"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetricsByPlugin contains an overview of the Agent metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AgentMetricsByPlugin(ctx context.Context, agentID string) (types.MetricsSummaryPlugin, error) {
	var out types.MetricsSummaryPlugin
	path := "/v2/agents/" + url.PathEscape(agentID) + "/metrics_by_plugin"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetricsOverTime contains Agent metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) AgentMetricsOverTime(ctx context.Context, agentID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v2/agents/" + url.PathEscape(agentID) + "/metrics_over_time?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetricsOverTimeByPlugin contains an overview of the Agent metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AgentMetricsOverTimeByPlugin(ctx context.Context, agentID string, params types.MetricsParams) (types.MetricsOverTimeByPlugin, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTimeByPlugin
	path := "/v2/agents/" + url.PathEscape(agentID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetrics contains an overview of the Core Instance metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) CoreInstanceMetrics(ctx context.Context, coreInstanceID string) (types.MetricsSummary, error) {
	var out types.MetricsSummary
	path := "/v2/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetricsByPlugin contains an overview of the Core Instance metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) CoreInstanceMetricsByPlugin(ctx context.Context, coreInstanceID string) (types.MetricsSummaryPlugin, error) {
	var out types.MetricsSummaryPlugin
	path := "/v2/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_by_plugin"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetricsOverTime contains Core Instance metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) CoreInstanceMetricsOverTime(ctx context.Context, coreInstanceID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v2/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_over_time?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetricsOverTimeByPlugin contains an overview of the Core Instance metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) CoreInstanceMetricsOverTimeByPlugin(ctx context.Context, coreInstanceID string, params types.MetricsParams) (types.MetricsOverTimeByPlugin, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTimeByPlugin
	path := "/v2/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetrics contains an overview of the Pipeline metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) PipelineMetrics(ctx context.Context, pipelineID string) (types.MetricsSummary, error) {
	var out types.MetricsSummary
	path := "/v2/pipelines/" + url.PathEscape(pipelineID) + "/metrics"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetricsByPlugin contains an overview of the Pipeline metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) PipelineMetricsByPlugin(ctx context.Context, pipelineID string) (types.MetricsSummaryPlugin, error) {
	var out types.MetricsSummaryPlugin
	path := "/v2/pipelines/" + url.PathEscape(pipelineID) + "/metrics_by_plugin"
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetricsOverTime contains Pipeline metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) PipelineMetricsOverTime(ctx context.Context, pipelineID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v2/pipelines/" + url.PathEscape(pipelineID) + "/metrics_over_time?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetricsOverTimeByPlugin contains an overview of the Pipeline metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) PipelineMetricsOverTimeByPlugin(ctx context.Context, pipelineID string, params types.MetricsParams) (types.MetricsOverTimeByPlugin, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTimeByPlugin
	path := "/v2/pipelines/" + url.PathEscape(pipelineID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
