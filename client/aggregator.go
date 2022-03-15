package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateAggregator within a project.
// The project in which the aggregator is created is parser from the authorization token.
// Users are not allowed to create aggregators.
func (c *Client) CreateAggregator(ctx context.Context, payload types.CreateAggregator) (types.CreatedAggregator, error) {
	var out types.CreatedAggregator
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators", payload, &out)
}

// Aggregators from a project in descending order.
func (c *Client) Aggregators(ctx context.Context, projectID string, params types.AggregatorsParams) (types.Aggregators, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out types.Aggregators
	path := "/v1/projects/" + url.PathEscape(projectID) + "/aggregators?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// Aggregator by ID.
func (c *Client) Aggregator(ctx context.Context, aggregatorID string) (types.Aggregator, error) {
	var out types.Aggregator
	return out, c.do(ctx, http.MethodGet, "/v1/aggregators/"+url.PathEscape(aggregatorID), nil, &out)
}

// UpdateAggregator by its ID.
func (c *Client) UpdateAggregator(ctx context.Context, aggregatorID string, payload types.UpdateAggregator) error {
	return c.do(ctx, http.MethodPatch, "/v1/aggregators/"+url.PathEscape(aggregatorID), payload, nil)
}

// DeleteAggregator by its ID.
func (c *Client) DeleteAggregator(ctx context.Context, aggregatorID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/aggregators/"+url.PathEscape(aggregatorID), nil, nil)
}
