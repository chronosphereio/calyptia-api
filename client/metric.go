package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// ProjectMetrics contains an overview of the aggregated metrics for a project.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) ProjectMetrics(ctx context.Context, projectID string, params types.MetricsParams) (types.ProjectMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.ProjectMetrics
	path := "/v1/projects/" + url.PathEscape(projectID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AgentMetrics contains an overview of the aggregated metrics for an agent.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AgentMetrics(ctx context.Context, agentID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.AgentMetrics
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineMetrics contains an overview of the aggregated metrics for a pipeline.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) PipelineMetrics(ctx context.Context, pipelineID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}
	var out types.AgentMetrics
	path := "/v1/pipeline_metrics/" + url.PathEscape(pipelineID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AggregatorMetrics contains an overview of the aggregated metrics for a project.
// It includes metrics link the amount of records, bytes, and errors per plugin.
func (c *Client) AggregatorMetrics(ctx context.Context, aggregatorID string, params types.MetricsParams) (types.AggregatorMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.AggregatorMetrics
	path := "/v1/aggregator_metrics/" + url.PathEscape(aggregatorID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// AggregatorPipelinesMetrics get the metrics for a set of pipelineIDs belonging to an aggregator (bulk mode).
func (c *Client) AggregatorPipelinesMetrics(ctx context.Context, aggregatorID string, params types.PipelinesMetricsParams) (types.PipelinesMetrics, error) {
	var out types.PipelinesMetrics
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	for _, pipelineID := range params.PipelineIDs {
		q.Add("pipeline_id", pipelineID)
	}

	return out, c.do(ctx, http.MethodGet, "/v1/aggregators/"+url.PathEscape(aggregatorID)+"/pipelines_metrics?"+q.Encode(), nil, &out)
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
		"start": []string{params.Start.String()},
	}

	var out types.MetricsOverTimeByPlugin
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/metrics_over_time_by_plugin?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
