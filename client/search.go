package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

func (c *Client) Search(ctx context.Context, in types.SearchQuery) ([]types.SearchResult, error) {
	var out []types.SearchResult

	q := url.Values{}
	q.Set("project_id", in.ProjectID)
	q.Set("resource", string(in.Resource))
	q.Set("term", in.Term)

	return out, c.do(ctx, http.MethodGet, "/v1/search?"+q.Encode(), nil, &out)
}
