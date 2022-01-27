package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

// ValidateConfig validates that an already parsed fluentbit or fluentd config is semantically valid.
// To parse the raw agent config take a look at https://github.com/calyptia/fluent-bit-config-parser.
func (c *Client) ValidateConfig(ctx context.Context, agentType types.AgentType, payload types.ValidatingConfig) (types.ValidatedConfig, error) {
	var out types.ValidatedConfig
	return out, c.do(ctx, http.MethodPost, "/v1/config_validate/"+url.PathEscape(string(agentType)), payload, &out)
}
