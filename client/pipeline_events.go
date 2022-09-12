package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// CreatePipeline within an aggregator.
// The pipeline name must be unique within the aggregator.
// The resource profile must exist already. If you don't provide one, it will default to "best-effort-low-resource".
// Use them to easily deploy configured agents to the aggregator.
func (c *Client) CreatePipelineEvents(ctx context.Context, aggregagotID string, payload types.CreatePipelineEvent) (types.CreatedPipelineEvent, error) {
	var out types.CreatedPipelineEvent
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(aggregagotID)+"/pipeline_event", payload, &out)
}
