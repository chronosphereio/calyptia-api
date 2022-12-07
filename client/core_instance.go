package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateCoreInstance within a project.
// The project in which the core instance is created is parser from the authorization token.
// Users are not allowed to create core instances.
func (c *Client) CreateCoreInstance(ctx context.Context, payload types.CreateCoreInstance) (types.CreatedCoreInstance, error) {
	var out types.CreatedCoreInstance
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators", payload, &out)
}

// CoreInstances from a project in descending order.
func (c *Client) CoreInstances(ctx context.Context, projectID string, params types.CoreInstancesParams) (types.CoreInstances, error) {
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
	if params.Tags != nil {
		q.Set("tags_query", *params.Tags)
	}
	if params.EnvironmentID != nil {
		q.Set("environment_id", *params.EnvironmentID)
	}

	var out types.CoreInstances
	path := "/v1/projects/" + url.PathEscape(projectID) + "/aggregators?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// CoreInstance by ID.
func (c *Client) CoreInstance(ctx context.Context, instanceID string) (types.CoreInstance, error) {
	var out types.CoreInstance
	return out, c.do(ctx, http.MethodGet, "/v1/aggregators/"+url.PathEscape(instanceID), nil, &out)
}

// UpdateCoreInstance by its ID.
func (c *Client) UpdateCoreInstance(ctx context.Context, instanceID string, payload types.UpdateCoreInstance) error {
	return c.do(ctx, http.MethodPatch, "/v1/aggregators/"+url.PathEscape(instanceID), payload, nil)
}

// DeleteCoreInstance by its ID.
func (c *Client) DeleteCoreInstance(ctx context.Context, instanceID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/aggregators/"+url.PathEscape(instanceID), nil, nil)
}

// CoreInstancePing by its ID.
func (c *Client) CoreInstancePing(ctx context.Context, instanceID string) (types.CoreInstancePingResponse, error) {
	var out types.CoreInstancePingResponse
	return out, c.do(ctx, http.MethodPost, "/v1/aggregators/"+url.PathEscape(instanceID)+"/ping", nil, &out)
}

// DeleteCoreInstances from a project passing a list of the IDs to be deleted.
func (c *Client) DeleteCoreInstances(ctx context.Context, projectID string, instanceIDs ...string) error {
	q := url.Values{}
	for _, id := range instanceIDs {
		q.Add("aggregator_id", id)
	}
	return c.do(ctx, http.MethodDelete, "/v1/projects/"+url.PathEscape(projectID)+"/aggregators?"+q.Encode(), nil, nil)
}
