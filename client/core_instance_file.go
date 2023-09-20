package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateCoreInstanceFile(ctx context.Context, in types.CreateCoreInstanceFile) (types.Created, error) {
	var out types.Created
	endpoint := "/v1/core_instances/" + url.PathEscape(in.CoreInstanceID) + "/files"
	return out, c.do(ctx, http.MethodPost, endpoint, in, &out)
}

func (c *Client) CoreInstanceFiles(ctx context.Context, in types.ListCoreInstanceFiles) (types.CoreInstanceFiles, error) {
	var out types.CoreInstanceFiles
	q := url.Values{}
	if in.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*in.Last), 10))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	endpoint := "/v1/core_instances/" + url.PathEscape(in.CoreInstanceID) + "/files"
	if len(q) != 0 {
		endpoint += "?" + q.Encode()
	}
	return out, c.do(ctx, http.MethodGet, endpoint, nil, &out)
}

// AllCoreInstanceFiles returns all core instance files for a given core instance
// by making multiple requests to the paginated API until all files have been
// retrieved.
func (c *Client) AllCoreInstanceFiles(ctx context.Context, coreInstanceID string) ([]types.CoreInstanceFile, error) {
	var out []types.CoreInstanceFile
	var cursor *string

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		page, err := c.CoreInstanceFiles(ctx, types.ListCoreInstanceFiles{
			CoreInstanceID: coreInstanceID,
			Before:         cursor,
		})
		if err != nil {
			return nil, err
		}

		out = append(out, page.Items...)
		if page.EndCursor == nil {
			break
		}

		cursor = page.EndCursor
	}

	return out, nil
}

func (c *Client) UpdateCoreInstanceFile(ctx context.Context, in types.UpdateCoreInstanceFile) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/core_instance_files/"+in.ID, in, &out)
}

func (c *Client) DeleteCoreInstanceFile(ctx context.Context, coreInstanceFileID string) (types.Deleted, error) {
	var out types.Deleted
	return out, c.do(ctx, http.MethodDelete, "/v1/core_instance_files/"+coreInstanceFileID, nil, &out)
}
