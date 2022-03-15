package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreatePipeline within an aggregator.
// The pipeline name must be unique within the aggregator.
// The resource profile must exist already. If you don't provide one, it will default to "best-effort-low-resource".
// Use them to easily deploy configured agents to the aggregator.
func (c *Client) CreatePipeline(ctx context.Context, aggregatorID string, payload types.CreatePipeline) (types.CreatedPipeline, error) {
	var out types.CreatedPipeline
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(aggregatorID)+"/pipelines", payload, &out)
}

// Pipelines from an aggregator in descending order.
func (c *Client) Pipelines(ctx context.Context, aggregatorID string, params types.PipelinesParams) (types.Pipelines, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out types.Pipelines
	path := "/v1/aggregators/" + url.PathEscape(aggregatorID) + "/pipelines?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// ProjectPipelines returns the entire set of pipelines from a project.
func (c *Client) ProjectPipelines(ctx context.Context, projectID string, params types.PipelinesParams) (types.Pipelines, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out types.Pipelines
	path := "/v1/projects/" + url.PathEscape(projectID) + "/aggregator_pipelines?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items)
}

// Pipeline by ID.
func (c *Client) Pipeline(ctx context.Context, pipelineID string) (types.Pipeline, error) {
	var out types.Pipeline
	return out, c.do(ctx, http.MethodGet, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID), nil, &out)
}

// UpdatePipeline by its ID.
func (c *Client) UpdatePipeline(ctx context.Context, pipelineID string, opts types.UpdatePipeline) (types.UpdatedPipeline, error) {
	var out types.UpdatedPipeline
	return out, c.do(ctx, http.MethodPatch, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID), opts, &out)
}

// DeletePipeline by its ID.
func (c *Client) DeletePipeline(ctx context.Context, pipelineID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID), nil, nil)
}
