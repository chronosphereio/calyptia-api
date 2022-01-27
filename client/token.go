package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateToken withint a project.
// These tokens are to authorize other applications to access the project.
// For example:
// - an agent might use it to register itself to the project.
// - you might create a new aggregator in the project using the aggregator CLI.
// - you might use it within the Calyptia CLI to grant access to your project.
func (c *Client) CreateToken(ctx context.Context, projectID string, payload types.CreateToken) (types.Token, error) {
	var out types.Token
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/tokens", payload, &out)
}

// Tokens from a project.
func (c *Client) Tokens(ctx context.Context, projectID string, params types.TokensParams) ([]types.Token, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, DefaultUintFormatBase))
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out []types.Token
	path := "/v1/projects/" + url.PathEscape(projectID) + "/tokens?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

// Token by ID.
func (c *Client) Token(ctx context.Context, tokenID string) (types.Token, error) {
	var out types.Token
	return out, c.do(ctx, http.MethodGet, "/v1/project_tokens/"+url.PathEscape(tokenID), nil, &out)
}

// UpdateToken by its ID.
func (c *Client) UpdateToken(ctx context.Context, tokenID string, opts types.UpdateToken) error {
	return c.do(ctx, http.MethodPatch, "/v1/project_tokens/"+url.PathEscape(tokenID), opts, nil)
}

// DeleteToken by its ID.
// Once deleted, any application that might has been using it, will stop working.
func (c *Client) DeleteToken(ctx context.Context, tokenID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/project_tokens/"+url.PathEscape(tokenID), nil, nil)
}
