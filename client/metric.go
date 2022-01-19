package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

func (c *Client) ProjectMetrics(ctx context.Context, projectID string, params types.MetricsParams) (types.ProjectMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.ProjectMetrics
	path := "/v1/projects/" + url.PathEscape(projectID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) AgentMetrics(ctx context.Context, agentID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.AgentMetrics
	path := "/v1/agents/" + url.PathEscape(agentID) + "/metrics?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) PipelineMetrics(ctx context.Context, pipelineID string, params types.MetricsParams) (types.AgentMetrics, error) {
	q := url.Values{
		"start":    []string{params.Start.String()},
		"interval": []string{params.Interval.String()},
	}

	var out types.AgentMetrics
	path := "/v1/pipeline_metrics/" + url.PathEscape(pipelineID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
