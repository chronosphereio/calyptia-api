package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// Members from a project in descending order.
func (c *Client) Members(ctx context.Context, projectID string, params types.MembersParams) ([]types.Membership, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}

	var out []types.Membership
	path := "/v1/projects/" + url.PathEscape(projectID) + "/members?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
