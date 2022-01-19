package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreatePipelineFile(ctx context.Context, pipelineID string, payload types.CreatePipelineFile) (types.CreatedPipelineFile, error) {
	var out types.CreatedPipelineFile
	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/files", payload, &out)
}

func (c *Client) PipelineFiles(ctx context.Context, pipelineID string, params types.PipelineFilesParams) ([]types.PipelineFile, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.PipelineFile
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/files?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) PipelineFile(ctx context.Context, fileID string) (types.PipelineFile, error) {
	var out types.PipelineFile
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_files/"+url.PathEscape(fileID), nil, &out)
}

func (c *Client) UpdatePipelineFile(ctx context.Context, fileID string, opts types.UpdatePipelineFile) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_files/"+url.PathEscape(fileID), opts, nil)
}

func (c *Client) DeletePipelineFile(ctx context.Context, fileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_files/"+url.PathEscape(fileID), nil, nil)
}
