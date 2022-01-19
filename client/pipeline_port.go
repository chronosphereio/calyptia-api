package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreatePipelinePort(ctx context.Context, pipelineID string, payload types.CreatePipelinePort) (types.CreatedPipelinePort, error) {
	var out types.CreatedPipelinePort
	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/ports", payload, &out)
}

func (c *Client) PipelinePorts(ctx context.Context, pipelineID string, params types.PipelinePortsParams) ([]types.PipelinePort, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.PipelinePort
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/ports?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) PipelinePort(ctx context.Context, portID string) (types.PipelinePort, error) {
	var out types.PipelinePort
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_ports/"+url.PathEscape(portID), nil, &out)
}

func (c *Client) UpdatePipelinePort(ctx context.Context, portID string, opts types.UpdatePipelinePort) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_ports/"+url.PathEscape(portID), opts, nil)
}

func (c *Client) DeletePipelinePort(ctx context.Context, portID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_ports/"+url.PathEscape(portID), nil, nil)
}
