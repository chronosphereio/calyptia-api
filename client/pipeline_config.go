package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) PipelineConfigHistory(ctx context.Context, pipelineID string, params types.PipelineConfigHistoryParams) ([]types.PipelineConfig, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.PipelineConfig
	path := "/v1/aggregator_pipelines/" + url.PathEscape(pipelineID) + "/config_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) PipelineConfig(ctx context.Context, configID string) (types.PipelineConfig, error) {
	var out types.PipelineConfig
	path := "/v1/pipeline_config_history/" + url.PathEscape(configID)
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
