package client

import (
	"context"
	"net/http"
)

// SendVerificationEmail sends a verification email to the user.
// When a new user registers within Calyptia Cloud with password, it should receive a verification email,
// in the case it didn't receive it, or the email expired already,
// use this endpoint to request a new one.
func (c *Client) SendVerificationEmail(ctx context.Context) error {
	return c.do(ctx, http.MethodGet, "/v1/verification_email", nil, nil)
}
