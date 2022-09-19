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

// CreatePipelineSecret within a pipeline.
// The given name is unique within the pipeline.
// These secrets can be referenced by their name within a fluentbit configuration file like so `{{secrets.thename}}`.
// Use them to hide sensible values from your config file.
func (c *Client) CreatePipelineSecret(ctx context.Context, pipelineID string, payload types.CreatePipelineSecret) (types.CreatedPipelineSecret, error) {
	var out types.CreatedPipelineSecret

	return out, c.do(ctx, http.MethodPost, "/v1/aggregator_pipelines/"+url.PathEscape(pipelineID)+"/secrets", payload, &out)
}

// PipelineSecrets in descending order.
func (c *Client) PipelineSecrets(ctx context.Context, pipelineID string, params types.PipelineSecretsParams) (types.PipelineSecrets, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.PipelineSecrets
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/secrets?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// PipelineSecret by ID.
func (c *Client) PipelineSecret(ctx context.Context, secretID string) (types.PipelineSecret, error) {
	var out types.PipelineSecret
	return out, c.do(ctx, http.MethodGet, "/v1/pipeline_secrets/"+url.PathEscape(secretID), nil, &out)
}

// UpdatePipelineSecret by its ID.
func (c *Client) UpdatePipelineSecret(ctx context.Context, secretID string, opts types.UpdatePipelineSecret) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipeline_secrets/"+url.PathEscape(secretID), opts, nil)
}

// DeletePipelineSecret by its ID.
// The secret cannot be deleted if some pipeline config is still referencing it;
// you must delete the pipeline first if you want to delete the secret.
func (c *Client) DeletePipelineSecret(ctx context.Context, secretID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/pipeline_secrets/"+url.PathEscape(secretID), nil, nil)
}
