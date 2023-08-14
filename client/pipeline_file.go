// Package client provides a client over the REST HTTP API of Calyptia Cloud.
//
//nolint:dupl // no need to remove duplication here.
package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreatePipelineFile within a pipeline.
// The given name is unique within the pipeline.
// These files can be referenced by their name within a fluentbit configuration file like so `{{files.thename}}`.
// Use them to share common stuff like parsers.
func (c *Client) CreatePipelineFile(ctx context.Context, pipelineID string, payload types.CreatePipelineFile) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/files", payload, &out)
}

// PipelineFiles in descending order.
func (c *Client) PipelineFiles(ctx context.Context, pipelineID string, params types.PipelineFilesParams) (types.PipelineFiles, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.PipelineFiles
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/files?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// PipelineFile by ID.
func (c *Client) PipelineFile(ctx context.Context, fileID string) (types.PipelineFile, error) {
	var out types.PipelineFile
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_files/"+url.PathEscape(fileID), nil, &out)
}

// UpdatePipelineFile by its ID.
func (c *Client) UpdatePipelineFile(ctx context.Context, fileID string, opts types.UpdatePipelineFile) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_files/"+url.PathEscape(fileID), opts, nil)
}

// DeletePipelineFile by its ID.
// The file cannot be deleted if some pipeline config is still referencing it;
// you must delete the pipeline first if you want to delete the file.
func (c *Client) DeletePipelineFile(ctx context.Context, fileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_files/"+url.PathEscape(fileID), nil, nil)
}
