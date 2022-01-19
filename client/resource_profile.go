package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateResourceProfile(ctx context.Context, aggregatorID string, payload types.CreateResourceProfile) (types.CreatedResourceProfile, error) {
	var out types.CreatedResourceProfile
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(aggregatorID)+"/resource_profiles", payload, &out)
}

func (c *Client) ResourceProfiles(ctx context.Context, aggregatorID string, params types.ResourceProfilesParams) ([]types.ResourceProfile, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.ResourceProfile
	path := "/v1/aggregators/" + url.PathEscape(aggregatorID) + "/resource_profiles?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) ResourceProfile(ctx context.Context, resourceProfileID string) (types.ResourceProfile, error) {
	var out types.ResourceProfile
	return out, c.do(ctx, http.MethodGet, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, &out)
}

func (c *Client) UpdateResourceProfile(ctx context.Context, resourceProfileID string, opts types.UpdateResourceProfile) error {
	return c.do(ctx, http.MethodPatch, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), opts, nil)
}

func (c *Client) DeleteResourceProfile(ctx context.Context, resourceProfileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, nil)
}
