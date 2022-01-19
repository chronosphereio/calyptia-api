package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateAggregator(ctx context.Context, payload types.CreateAggregator) (types.CreatedAggregator, error) {
	var out types.CreatedAggregator
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators", payload, &out)
}

func (c *Client) Aggregators(ctx context.Context, projectID string, params types.AggregatorsParams) ([]types.Aggregator, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out []types.Aggregator
	path := "/v1/projects/" + url.PathEscape(projectID) + "/aggregators?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) Aggregator(ctx context.Context, aggregatorID string) (types.Aggregator, error) {
	var out types.Aggregator
	return out, c.do(ctx, http.MethodGet, "/v1/aggregators/"+url.PathEscape(aggregatorID), nil, &out)
}

func (c *Client) UpdateAggregator(ctx context.Context, aggregatorID string, payload types.UpdateAggregator) error {
	return c.do(ctx, http.MethodPatch, "/v1/aggregators/"+url.PathEscape(aggregatorID), payload, nil)
}

func (c *Client) DeleteAggregator(ctx context.Context, aggregatorID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/aggregators/"+url.PathEscape(aggregatorID), nil, nil)
}
