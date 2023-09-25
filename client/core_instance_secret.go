package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateCoreInstanceSecret(ctx context.Context, in types.CreateCoreInstanceSecret) (types.Created, error) {
	var out types.Created
	endpoint := "/v1/core_instances/" + url.PathEscape(in.CoreInstanceID) + "/secrets"
	return out, c.do(ctx, http.MethodPost, endpoint, in, &out)
}

func (c *Client) CoreInstanceSecrets(ctx context.Context, in types.ListCoreInstanceSecrets) (types.CoreInstanceSecrets, error) {
	var out types.CoreInstanceSecrets
	q := url.Values{}
	if in.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*in.Last), 10))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	endpoint := "/v1/core_instances/" + url.PathEscape(in.CoreInstanceID) + "/secrets"
	if len(q) != 0 {
		endpoint += "?" + q.Encode()
	}
	return out, c.do(ctx, http.MethodGet, endpoint, nil, &out)
}

// AllCoreInstanceSecrets returns all core instance secrets for a given core instance
// by making multiple requests to the paginated API until all secrets have been
// retrieved.
func (c *Client) AllCoreInstanceSecrets(ctx context.Context, coreInstanceID string) ([]types.CoreInstanceSecret, error) {
	var out []types.CoreInstanceSecret
	var cursor *string

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		page, err := c.CoreInstanceSecrets(ctx, types.ListCoreInstanceSecrets{
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

func (c *Client) UpdateCoreInstanceSecret(ctx context.Context, in types.UpdateCoreInstanceSecret) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/core_instance_secrets/"+in.ID, in, &out)
}

func (c *Client) DeleteCoreInstanceSecret(ctx context.Context, coreInstanceSecretID string) (types.Deleted, error) {
	var out types.Deleted
	return out, c.do(ctx, http.MethodDelete, "/v1/core_instance_secrets/"+coreInstanceSecretID, nil, &out)
}
