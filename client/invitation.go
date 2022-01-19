package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateInvitation(ctx context.Context, projectID string, payload types.CreateInvitation) error {
	return c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/invite", payload, nil)
}

func (c *Client) AcceptInvitation(ctx context.Context, payload types.AcceptInvitation) error {
	return c.do(ctx, http.MethodPost, "/v1/join_project", payload, nil)
}
