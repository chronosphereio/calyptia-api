package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateEnvironment within a project.
func (c *Client) CreateEnvironment(ctx context.Context, projectID string, payload types.CreateEnvironment) (types.CreatedEnvironment, error) {
	var out types.CreatedEnvironment
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/environments", payload, &out)
}

// Environments from the given project in descending order.
func (c *Client) Environments(ctx context.Context, projectID string, params types.EnvironmentsParams) (types.Environments, error) {
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

	var out types.Environments
	path := "/v1/projects/" + url.PathEscape(projectID) + "/environments?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// UpdateEnvironment by its ID.
func (c *Client) UpdateEnvironment(ctx context.Context, environmentID string, payload types.UpdateEnvironment) error {
	return c.do(ctx, http.MethodPatch, "/v1/environments/"+url.PathEscape(environmentID), payload, nil)
}

// DeleteEnvironment by its ID.
func (c *Client) DeleteEnvironment(ctx context.Context, environmentID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/environments/"+url.PathEscape(environmentID), nil, nil)
}
