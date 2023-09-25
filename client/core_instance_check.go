package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateCoreInstanceCheck within a core_instance.
func (c *Client) CreateCoreInstanceCheck(ctx context.Context, coreInstanceID string, payload types.CreateCoreInstanceCheck) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/core_instances/"+url.PathEscape(coreInstanceID)+"/checks", payload, &out)
}

// CoreInstanceChecks in descending order.
func (c *Client) CoreInstanceChecks(ctx context.Context, coreInstanceID string, params types.CoreInstanceChecksParams) (types.CoreInstanceChecks, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.CoreInstanceChecks
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/checks?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// CoreInstanceCheck by ID.
func (c *Client) CoreInstanceCheck(ctx context.Context, checkID string) (types.CoreInstanceCheck, error) {
	var out types.CoreInstanceCheck
	return out, c.do(ctx, http.MethodGet, "/v1/core_instance_checks/"+url.PathEscape(checkID), nil, &out)
}

// UpdateCoreInstanceCheck by its ID.
func (c *Client) UpdateCoreInstanceCheck(ctx context.Context, checkID string, opts types.UpdateCoreInstanceCheck) error {
	return c.do(ctx, http.MethodPatch, "/v1/core_instance_checks/"+url.PathEscape(checkID), opts, nil)
}

// DeleteCoreInstanceCheck by its ID.
func (c *Client) DeleteCoreInstanceCheck(ctx context.Context, checkID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/core_instance_checks/"+url.PathEscape(checkID), nil, nil)
}
