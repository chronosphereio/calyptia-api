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
func (c *Client) PipelineConfigHistory(ctx context.Context, pipelineID string, params types.PipelineConfigHistoryParams) ([]types.PipelineConfig, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, defaultUintFormatBase))
	}

	var out []types.PipelineConfig
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/config_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
