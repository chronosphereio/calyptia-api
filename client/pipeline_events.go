package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// CreatePipelineEvent for storing extended statuses.
func (c *Client) CreatePipelineEvent(ctx context.Context, pipelineID string, payload types.CreatePipelineEvent) (types.CreatedPipelineEvent, error) {
	var out types.CreatedPipelineEvent
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/events", payload, &out)
}
