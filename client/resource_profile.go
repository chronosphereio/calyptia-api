// Package client provides a client over the REST HTTP API of Calyptia Cloud.
//
//nolint:dupl // no need to remove duplication here.
package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateResourceProfile within a core instance.
// A resource profile is a specification of a resource used during the deployment of a pipeline.
// By default, when you setup a core instance, Calyptia Cloud will generate 3 resource profiles for you:
// - high-performance-guaranteed-delivery.
// - high-performance-optimal-throughput.
// - best-effort-low-resource.
func (c *Client) CreateResourceProfile(ctx context.Context, instanceID string, payload types.CreateResourceProfile) (types.CreatedResourceProfile, error) {
	var out types.CreatedResourceProfile
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(instanceID)+"/resource_profiles", payload, &out)
}

// ResourceProfiles from a core instance in descending order.
func (c *Client) ResourceProfiles(ctx context.Context, instanceID string, params types.ResourceProfilesParams) (types.ResourceProfiles, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.ResourceProfiles
	path := "/v1/aggregators/" + url.PathEscape(instanceID) + "/resource_profiles?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// ResourceProfile by ID.
func (c *Client) ResourceProfile(ctx context.Context, resourceProfileID string) (types.ResourceProfile, error) {
	var out types.ResourceProfile
	return out, c.do(ctx, http.MethodGet, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, &out)
}

// UpdateResourceProfile by its ID.
func (c *Client) UpdateResourceProfile(ctx context.Context, resourceProfileID string, opts types.UpdateResourceProfile) error {
	return c.do(ctx, http.MethodPatch, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), opts, nil)
}

// DeleteResourceProfile by its ID.
// The profile cannot be deleted if some pipeline is still referencing it;
// you must delete the pipeline first if you want to delete the profile.
func (c *Client) DeleteResourceProfile(ctx context.Context, resourceProfileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/resource_profiles/"+url.PathEscape(resourceProfileID), nil, nil)
}
