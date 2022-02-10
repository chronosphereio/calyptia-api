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

// ValidateConfigV2 validates that an already parsed fluentbit(only) config to check if semantically valid
// To parse the raw agent config take a look at https://github.com/calyptia/fluent-bit-config-parser.
func (c *Client) ValidateConfigV2(ctx context.Context, payload types.ValidatingConfig) (types.ValidatedConfigV2, error) {
	var out types.ValidatedConfigV2
	return out, c.do(ctx, http.MethodPost, "/v1/config_validate_v2", payload, &out)
}
