package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateFleetFile within a fleet.
// The given name is unique within the fleet.
// These files can be referenced by their name within a fluentbit configuration file like so `{{files.thename}}`.
// Use them to share common stuff like parsers.
func (c *Client) CreateFleetFile(ctx context.Context, fleetID string, payload types.CreateFleetFile) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/fleets/"+url.PathEscape(fleetID)+"/files", payload, &out)
}

// FleetFiles in descending order.
func (c *Client) FleetFiles(ctx context.Context, fleetID string, params types.FleetFilesParams) (types.FleetFiles, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.FleetFiles
	path := "/v1/fleets/" + url.PathEscape(fleetID) + "/files?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// FleetFile by ID.
func (c *Client) FleetFile(ctx context.Context, fileID string) (types.FleetFile, error) {
	var out types.FleetFile
	return out, c.do(ctx, http.MethodGet, "/v1/fleet_files/"+url.PathEscape(fileID), nil, &out)
}

// UpdateFleetFile by its ID.
func (c *Client) UpdateFleetFile(ctx context.Context, fileID string, opts types.UpdateFleetFile) error {
	return c.do(ctx, http.MethodPatch, "/v1/fleet_files/"+url.PathEscape(fileID), opts, nil)
}

// DeleteFleetFile by its ID.
// The file cannot be deleted if some fleet config is still referencing it;
// you must delete the fleet first if you want to delete the file.
func (c *Client) DeleteFleetFile(ctx context.Context, fileID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/fleet_files/"+url.PathEscape(fileID), nil, nil)
}
