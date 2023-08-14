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

// CreateIngestCheck within a core_instance.
func (c *Client) CreateIngestCheck(ctx context.Context, coreInstanceID string, payload types.CreateIngestCheck) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/core_instances/"+url.PathEscape(coreInstanceID)+"/ingest_checks", payload, &out)
}

// IngestChecks in descending order.
func (c *Client) IngestChecks(ctx context.Context, coreInstanceID string, params types.IngestChecksParams) (types.IngestChecks, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.IngestChecks
	path := "/v1/core_instances/" + url.PathEscape(coreInstanceID) + "/ingest_checks?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// IngestCheck by ID.
func (c *Client) IngestCheck(ctx context.Context, checkID string) (types.IngestCheck, error) {
	var out types.IngestCheck
	return out, c.do(ctx, http.MethodGet, "/v1/ingest_checks/"+url.PathEscape(checkID), nil, &out)
}

// UpdateIngestCheck by its ID.
func (c *Client) UpdateIngestCheck(ctx context.Context, checkID string, opts types.UpdateIngestCheck) error {
	return c.do(ctx, http.MethodPatch, "/v1/ingest_checks/"+url.PathEscape(checkID), opts, nil)
}

// DeleteIngestCheck by its ID.
func (c *Client) DeleteIngestCheck(ctx context.Context, checkID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/ingest_checks/"+url.PathEscape(checkID), nil, nil)
}
