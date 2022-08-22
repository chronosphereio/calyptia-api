package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// PipelineStatusHistory in descending order.
// Every time a pipeline status is changed, a new history entry with the change is created.
func (c *Client) PipelineStatusHistory(ctx context.Context, pipelineID string, params types.PipelineStatusHistoryParams) (types.PipelineStatusHistory, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Status != nil {
		q.Set("status", string(*params.Status))
	}
	if params.ConfigFormat != nil {
		q.Set("config_format", string(*params.ConfigFormat))
	}

	var out types.PipelineStatusHistory
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/status_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}
