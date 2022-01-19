package client

import (
	"context"
	"net/http"
)

func (c *Client) SendVerificationEmail(ctx context.Context) error {
	return c.do(ctx, http.MethodPost, "/v1/verification_email", nil, nil)
}
