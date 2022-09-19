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

// CreatePipelinePort within a pipeline.
// Ports can automatically be parsed from a config file, but this action allows you to programatically add more.
func (c *Client) CreatePipelinePort(ctx context.Context, pipelineID string, payload types.CreatePipelinePort) (types.CreatedPipelinePort, error) {
	var out types.CreatedPipelinePort
	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/ports", payload, &out)
}

// PipelinePorts in descending order.
func (c *Client) PipelinePorts(ctx context.Context, pipelineID string, params types.PipelinePortsParams) (types.PipelinePorts, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.PipelinePorts
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/ports?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// PipelinePort by ID.
func (c *Client) PipelinePort(ctx context.Context, portID string) (types.PipelinePort, error) {
	var out types.PipelinePort
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_ports/"+url.PathEscape(portID), nil, &out)
}

// UpdatePipelinePort by its ID.
func (c *Client) UpdatePipelinePort(ctx context.Context, portID string, opts types.UpdatePipelinePort) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_ports/"+url.PathEscape(portID), opts, nil)
}

// DeletePipelinePort by its ID.
func (c *Client) DeletePipelinePort(ctx context.Context, portID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_ports/"+url.PathEscape(portID), nil, nil)
}
