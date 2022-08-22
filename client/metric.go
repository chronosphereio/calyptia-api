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
