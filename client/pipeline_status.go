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
func (c *Client) PipelineStatusHistory(ctx context.Context, pipelineID string, params types.PipelineStatusHistoryParams) ([]types.PipelineStatus, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, defaultUintFormatBase))
	}

	var out []types.PipelineStatus
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/status_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
