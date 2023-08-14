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

// CreatePipelineCheck within a pipeline.
func (c *Client) CreatePipelineCheck(ctx context.Context, pipelineID string, payload types.CreatePipelineCheck) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/checks", payload, &out)
}

// PipelineChecks in descending order.
func (c *Client) PipelineChecks(ctx context.Context, pipelineID string, params types.PipelineChecksParams) (types.PipelineChecks, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.PipelineChecks
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/checks?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// PipelineCheck by ID.
func (c *Client) PipelineCheck(ctx context.Context, checkID string) (types.PipelineCheck, error) {
	var out types.PipelineCheck
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_checks/"+url.PathEscape(checkID), nil, &out)
}

// UpdatePipelineCheck by its ID.
func (c *Client) UpdatePipelineCheck(ctx context.Context, checkID string, opts types.UpdatePipelineCheck) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_checks/"+url.PathEscape(checkID), opts, nil)
}

// DeletePipelineCheck by its ID.
func (c *Client) DeletePipelineCheck(ctx context.Context, checkID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_checks/"+url.PathEscape(checkID), nil, nil)
}
