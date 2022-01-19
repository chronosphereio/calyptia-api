package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// CreateInvitation to a project.
// This will send an invitation email with a link to join to the email address provided.
func (c *Client) CreateInvitation(ctx context.Context, projectID string, payload types.CreateInvitation) error {
	return c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/invite", payload, nil)
}

// AcceptInvitation to a project.
// The project to which you join is parsed from the token.
func (c *Client) AcceptInvitation(ctx context.Context, payload types.AcceptInvitation) error {
	return c.do(ctx, http.MethodPost, "/v1/join_project", payload, nil)
}
