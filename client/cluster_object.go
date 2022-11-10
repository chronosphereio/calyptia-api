// Package client provides a client over the REST HTTP API of Calyptia Cloud.
package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateClusterObject within a core_instance.
func (c *Client) CreateClusterObject(ctx context.Context, coreInstanceID string, payload types.CreateClusterObject) (types.CreatedClusterObject, error) {
	var out types.CreatedClusterObject
	return out, c.do(ctx, http.MethodPost, "/v1/core_instances/"+url.PathEscape(coreInstanceID)+"/cluster_objects", payload, &out)
}

// ClusterObjects in descending order.
func (c *Client) ClusterObjects(ctx context.Context, coreInstanceID string, params types.ClusterObjectParams) (types.ClusterObjects, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}
	if params.Kind != nil {
		q.Set("kind", string(*params.Kind))
	}

	var out types.ClusterObjects
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/cluster_objects?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// ClusterObject by ID.
func (c *Client) ClusterObject(ctx context.Context, checkID string) (types.ClusterObject, error) {
	var out types.ClusterObject
	return out, c.do(ctx, http.MethodGet, "/v1/cluster_objects/"+url.PathEscape(checkID), nil, &out)
}

// UpdateClusterObject by its ID.
func (c *Client) UpdateClusterObject(ctx context.Context, checkID string, opts types.UpdateClusterObject) error {
	return c.do(ctx, http.MethodPatch, "/v1/cluster_objects/"+url.PathEscape(checkID), opts, nil)
}

// DeleteClusterObject by its ID.
func (c *Client) DeleteClusterObject(ctx context.Context, checkID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/cluster_objects/"+url.PathEscape(checkID), nil, nil)
}
