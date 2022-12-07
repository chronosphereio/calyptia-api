package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// ProjectMetricsV1 contains an overview of the aggregated metrics for a project.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) ProjectMetricsV1(ctx context.Context, projectID string, params types.MetricsParams) (types.ProjectMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.ProjectMetrics
	path := "/v1/projects/" + url.PathEscape(projectID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetricsV1 contains an overview of the aggregated metrics for an agent.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AgentMetricsV1(ctx context.Context, agentID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.AgentMetrics
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetricsV1 contains an overview of the aggregated metrics for a pipeline.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) PipelineMetricsV1(ctx context.Context, pipelineID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}
	var out types.AgentMetrics
	path := "/v1/pipeline_metrics/" + url.PathEscape(pipelineID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetricsV1 contains an overview of the aggregated metrics for a project.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) CoreInstanceMetricsV1(ctx context.Context, instanceID string, params types.MetricsParams) (types.CoreInstanceMetricsV1, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.CoreInstanceMetricsV1
	path := "/v1/aggregator_metrics/" + url.PathEscape(instanceID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelinesMetricsV1 get the metrics for a set of pipelineIDs belonging to a core instance (bulk mode).
func (c *Client) PipelinesMetricsV1(ctx context.Context, instanceID string, params types.PipelinesMetricsParams) (types.PipelinesMetrics, error) {
	var out types.PipelinesMetrics
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	for _, pipelineID := range params.PipelineIDs {
		q.Add("pipeline_id", pipelineID)
	}

	return out, c.do(ctx, http.MethodGet, "/v1/aggregators/"+url.PathEscape(instanceID)+"/pipelines_metrics?"+q.Encode(), nil, &out)
}

// CoreInstanceMetrics contains an overview of the Core Instance metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) CoreInstanceMetrics(ctx context.Context, coreInstanceID string, params types.MetricsParams) (types.MetricsSummary, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummary
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceMetricsByPlugin contains an overview of the Core Instance metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) CoreInstanceMetricsByPlugin(ctx context.Context, coreInstanceID string, params types.MetricsParams) (types.MetricsSummaryPlugin, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummaryPlugin
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// CoreInstanceOverTimeMetrics contains Core Instance metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) CoreInstanceOverTimeMetrics(ctx context.Context, coreInstanceID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_over_time?" + q.Encode()
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
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetrics contains an overview of the Pipeline metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) PipelineMetrics(ctx context.Context, pipelineID string, params types.MetricsParams) (types.MetricsSummary, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummary
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetricsByPlugin contains an overview of the Pipeline metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) PipelineMetricsByPlugin(ctx context.Context, pipelineID string, params types.MetricsParams) (types.MetricsSummaryPlugin, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummaryPlugin
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/metrics_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineOverTimeMetrics contains Pipeline metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) PipelineOverTimeMetrics(ctx context.Context, pipelineID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/metrics_over_time?" + q.Encode()
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
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetrics contains an overview of the Agent metrics.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) AgentMetrics(ctx context.Context, agentID string, params types.MetricsParams) (types.MetricsSummary, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummary
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetricsByPlugin contains an overview of the Agent metrics.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AgentMetricsByPlugin(ctx context.Context, agentID string, params types.MetricsParams) (types.MetricsSummaryPlugin, error) {
	q := url.Values{
		"start": []string{params.Start.String()},
	}

	var out types.MetricsSummaryPlugin
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentOverTimeMetrics contains Agent metrics overtime.
// It includes metrics link the amount of records, bytes, and errors.
func (c *Client) AgentOverTimeMetrics(ctx context.Context, agentID string, params types.MetricsParams) (types.MetricsOverTime, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.MetricsOverTime
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics_over_time?" + q.Encode()
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
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
