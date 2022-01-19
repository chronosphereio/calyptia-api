package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

func (c *Client) ValidateConfig(ctx context.Context, agentType types.AgentType, payload types.ValidatingConfig) (types.ValidatedConfig, error) {
	var out types.ValidatedConfig
	return out, c.do(ctx, http.MethodPost, "/v1/config_validate/"+url.PathEscape(string(agentType)), payload, &out)
}
