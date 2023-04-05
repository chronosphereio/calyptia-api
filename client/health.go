package client

import (
	"context"
	"net/http"

	"github.com/calyptia/api/types"
)

func (c *Client) Health(ctx context.Context) (types.Health, error) {
	var out types.Health
	return out, c.do(ctx, http.MethodGet, "/healthz", nil, &out)
}
