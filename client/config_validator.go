package client

import (
	"context"
	"net/http"

	"github.com/calyptia/api/types"
)

func (c *Client) ValidateFluentbitConfig(ctx context.Context, in types.ValidateFluentbitConfig) error {
	return c.do(ctx, http.MethodPost, "/v2/validate_fluentbit_config", in, nil)
}
