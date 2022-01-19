package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreatePipelineSecret(ctx context.Context, pipelineID string, payload types.CreatePipelineSecret) (types.CreatedPipelineSecret, error) {
	var out types.CreatedPipelineSecret

	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/secrets", payload, &out)
}

func (c *Client) PipelineSecrets(ctx context.Context, pipelineID string, params types.PipelineSecretsParams) ([]types.PipelineSecret, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.PipelineSecret
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/secrets?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) PipelineSecret(ctx context.Context, secretID string) (types.PipelineSecret, error) {
	var out types.PipelineSecret
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_secrets/"+url.PathEscape(secretID), nil, &out)
}

func (c *Client) UpdatePipelineSecret(ctx context.Context, secretID string, opts types.UpdatePipelineSecret) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_secrets/"+url.PathEscape(secretID), opts, nil)
}

func (c *Client) DeletePipelineSecret(ctx context.Context, secretID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_secrets/"+url.PathEscape(secretID), nil, nil)
}
