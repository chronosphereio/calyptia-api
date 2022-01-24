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
		q.Set("last", strconv.FormatUint(*params.Last, DefaultUintFormatBase))
	}

	var out []types.PipelineConfig
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/config_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// PipelineConfig retrieves a single pipeline config history entry by its ID.
func (c *Client) PipelineConfig(ctx context.Context, configID string) (types.PipelineConfig, error) {
	var out types.PipelineConfig
	path := "/v1/pipeline_config_history/" + url.PathEscape(configID)
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
