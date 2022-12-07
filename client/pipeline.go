package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreatePipeline within a core instance.
// The pipeline name must be unique within the core instance.
// The resource profile must exist already. If you don't provide one, it will default to "best-effort-low-resource".
// Use them to easily deploy configured agents to the core instance.
func (c *Client) CreatePipeline(ctx context.Context, instanceID string, payload types.CreatePipeline) (types.CreatedPipeline, error) {
	var out types.CreatedPipeline
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(instanceID)+"/pipelines", payload, &out)
}

// Pipelines from a core instance in descending order.
func (c *Client) Pipelines(ctx context.Context, instanceID string, params types.PipelinesParams) (types.Pipelines, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}
	if params.Tags != nil {
		q.Set("tags_query", *params.Tags)
	}
	if params.ConfigFormat != nil {
		q.Set("config_format", string(*params.ConfigFormat))
	}
	if params.RenderWithConfigSections {
		q.Set("render_with_config_sections", "true")
	}

	var out types.Pipelines
	path := "/v1/aggregators/" + url.PathEscape(instanceID) + "/pipelines?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// ProjectPipelines returns the entire set of pipelines from a project.
func (c *Client) ProjectPipelines(ctx context.Context, projectID string, params types.PipelinesParams) (types.Pipelines, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}
	if params.Tags != nil {
		q.Set("tags_query", *params.Tags)
	}
	if params.ConfigFormat != nil {
		q.Set("config_format", string(*params.ConfigFormat))
	}
	if params.RenderWithConfigSections {
		q.Set("render_with_config_sections", "true")
	}

	var out types.Pipelines
	path := "/v1/projects/" + url.PathEscape(projectID) + "/aggregator_pipelines?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// Pipeline by ID.
func (c *Client) Pipeline(ctx context.Context, pipelineID string, params types.PipelineParams) (types.Pipeline, error) {
	q := url.Values{}
	if params.ConfigFormat != nil {
		q.Set("config_format", string(*params.ConfigFormat))
	}
	if params.RenderWithConfigSections {
		q.Set("render_with_config_sections", "true")
	}

	var out types.Pipeline
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
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

// DeletePipelines from a core instance passing a list of the IDs to be deleted.
func (c *Client) DeletePipelines(ctx context.Context, instanceID string, pipelineIDs ...string) error {
	q := url.Values{}
	for _, id := range pipelineIDs {
		q.Add("pipeline_id", id)
	}
	return c.do(ctx, http.MethodDelete, "/v1/aggregators/"+url.PathEscape(instanceID)+"/pipelines?"+q.Encode(), nil, nil)
}

// UpdatePipelineClusterObjects update a list of related cluster objects to a pipeline.
func (c *Client) UpdatePipelineClusterObjects(ctx context.Context, pipelineID string, opts types.UpdatePipelineClusterObjects) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/cluster_objects", opts, nil)
}

// DeletePipelineClusterObjects un-relate a list of cluster objects from a pipeline.
func (c *Client) DeletePipelineClusterObjects(ctx context.Context, pipelineID string, clusterObjectIDs ...string) error {
	q := url.Values{}
	for _, id := range clusterObjectIDs {
		q.Add("cluster_object_id", id)
	}
	return c.do(ctx, http.MethodDelete, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/cluster_objects?"+q.Encode(), nil, nil)
}

// PipelineClusterObjects returns the entire set of cluster objects associated to a pipeline.
func (c *Client) PipelineClusterObjects(ctx context.Context, pipelineID string, params types.PipelineClusterObjectsParams) (types.ClusterObjects, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.ClusterObjects
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/cluster_objects?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}
