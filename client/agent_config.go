package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// AgentConfigHistory in descending order.
// Every time an agent config is updated, a new history entry with the change is created.
func (c *Client) AgentConfigHistory(ctx context.Context, agentID string, params types.AgentConfigHistoryParams) (types.AgentConfigHistory, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.AgentConfigHistory
	path := "/v1/agents/" + url.PathEscape(agentID) + "/config_history?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items)
}
