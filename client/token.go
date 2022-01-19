package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateToken(ctx context.Context, projectID string, payload types.CreateToken) (types.Token, error) {
	var out types.Token
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/aggregators", payload, &out)
}

func (c *Client) Tokens(ctx context.Context, projectID string, params types.TokensParams) ([]types.Token, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out []types.Token
	path := "/v1/projects/" + url.PathEscape(projectID) + "/tokens?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) Token(ctx context.Context, tokenID string) (types.Token, error) {
	var out types.Token
	return out, c.do(ctx, http.MethodGet, "/v1/project_tokens/"+url.PathEscape(tokenID), nil, &out)
}

func (c *Client) UpdateToken(ctx context.Context, tokenID string, opts types.UpdateToken) error {
	return c.do(ctx, http.MethodPatch, "/v1/project_tokens/"+url.PathEscape(tokenID), opts, nil)
}

func (c *Client) DeleteToken(ctx context.Context, tokenID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/project_tokens/"+url.PathEscape(tokenID), nil, nil)
}
