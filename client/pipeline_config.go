package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// PipelineConfigHistory in descending order.
// Every time a pipeline config is updated, a new history entry with the change is created.
func (c *Client) PipelineConfigHistory(ctx context.Context, pipelineID string, params types.PipelineConfigHistoryParams) (types.PipelineConfigHistory, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.ConfigFormat != nil {
		q.Set("config_format", string(*params.ConfigFormat))
	}

	var out types.PipelineConfigHistory
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/config_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}
