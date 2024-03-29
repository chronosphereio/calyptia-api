package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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
func (c *Client) Pipelines(ctx context.Context, in types.PipelinesParams) (types.Pipelines, error) {
	q := url.Values{}
	if in.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*in.Last), uintBase))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	if in.Kind != nil {
		q.Set("kind", string(*in.Kind))
	}
	if in.Name != nil {
		q.Set("name", *in.Name)
	}
	if in.TagsQuery != nil {
		q.Set("tags_query", *in.TagsQuery)
	}
	if in.ConfigFormat != nil {
		q.Set("config_format", string(*in.ConfigFormat))
	}
	if in.RenderWithConfigSections {
		q.Set("render_with_config_sections", "true")
	}

	if in.IncludeObjects != nil {
		var objects []string
		if in.IncludeObjects.Ports {
			objects = append(objects, "ports")
		}
		if in.IncludeObjects.Files {
			objects = append(objects, "files")
		}
		if in.IncludeObjects.Secrets {
			objects = append(objects, "secrets")
		}
		if len(objects) > 0 {
			q.Set("include", strings.Join(objects, ","))
		}
	}
	var out types.Pipelines
	var path string
	switch {
	case in.ProjectID != nil:
		path = "/v1/projects/" + url.PathEscape(*in.ProjectID) + "/pipelines?" + q.Encode()
	case in.CoreInstanceID != nil:
		path = "/v1/core_instances/" + url.PathEscape(*in.CoreInstanceID) + "/pipelines?" + q.Encode()
	default:
		return out, errors.New("either project_id or core_instance_id must be provided")
	}

	return out, c.do(ctx, http.MethodGet, path, nil, &out)
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
	var errs error
	for _, chunk := range makeChunks(pipelineIDs, chunksSizeByID) {
		q := url.Values{}
		for _, id := range chunk {
			q.Add("pipeline_id", id)
		}
		err := c.do(ctx, http.MethodDelete, "/v1/aggregators/"+url.PathEscape(instanceID)+"/pipelines?"+q.Encode(), nil, nil)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}

// UpdatePipelineClusterObjects update a list of related cluster objects to a pipeline.
func (c *Client) UpdatePipelineClusterObjects(ctx context.Context, pipelineID string, opts types.UpdatePipelineClusterObjects) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/cluster_objects", opts, nil)
}

// DeletePipelineClusterObjects un-relate a list of cluster objects from a pipeline.
func (c *Client) DeletePipelineClusterObjects(ctx context.Context, pipelineID string, clusterObjectIDs ...string) error {
	var errs error
	for _, chunk := range makeChunks(clusterObjectIDs, chunksSizeByID) {
		q := url.Values{}
		for _, id := range chunk {
			q.Add("cluster_object_id", id)
		}

		err := c.do(ctx, http.MethodDelete, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/cluster_objects?"+q.Encode(), nil, nil)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
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

// UpdatePipelineMetadata modifies the metadata associated with a specific pipeline.
//
// Parameters:
// - `ctx`: The context for executing the HTTP request.
// - `pipelineID`: The unique identifier for the target pipeline.
// - `opts`: A struct containing the key-value pair(s) to update or add to the pipeline's metadata.
//
// Returns:
// - An error if the request fails, if the pipeline with the given ID is not found, or if the provided metadata is invalid.
//
// The function constructs a request URL using the provided pipelineID and then sends a PATCH request
// to update the pipeline's metadata with the key-value pairs specified in the `opts` parameter.
// If a metadata key already exists, its value will be updated; if the key does not exist, it will be added.
func (c *Client) UpdatePipelineMetadata(ctx context.Context, pipelineID string, opts types.UpdatePipelineMetadata) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/metadata", opts, nil)
}

// PipelineMetadata retrieves the metadata associated with a given pipeline.
//
// Parameters:
//   - `ctx`: The context for executing the HTTP request.
//   - `pipelineID`: The unique identifier for the desired pipeline.
//   - `keys`: An optional list of specific metadata keys to retrieve. If no keys are provided,
//     all metadata for the pipeline will be returned.
//
// Returns:
// - A `types.PipelineMetadata` containing the requested metadata.
// - An error if the request fails or if the pipeline with the given ID is not found.
//
// The function constructs a request URL by appending the provided pipelineID and any specified keys
// as query parameters. It then performs a GET request to fetch the metadata.
func (c *Client) PipelineMetadata(ctx context.Context, pipelineID string, keys ...string) (types.PipelineMetadata, error) {
	q := url.Values{}
	for _, key := range keys {
		q.Add("key", key)
	}
	var out types.PipelineMetadata
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/metadata?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
