package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateAgentError(ctx context.Context, in types.CreateAgentError) (types.CreatedAgentError, error) {
	var out types.CreatedAgentError
	return out, c.do(ctx, http.MethodPost, "/v1/agents/"+in.AgentID+"/errors", in, &out)
}

func (c *Client) AgentErrors(ctx context.Context, in types.ListAgentErrors) (types.AgentErrors, error) {
	var out types.AgentErrors

	if in.AgentID == nil && in.FleetID == nil {
		return out, errors.New("either agent ID or fleet ID required")
	}

	var path string
	if in.AgentID != nil {
		path = "/v1/agents/" + *in.AgentID + "/errors"
	}
	if in.FleetID != nil {
		path = "/v1/fleets/" + *in.FleetID + "/agent_errors"
	}

	q := url.Values{}
	if in.Dismissed != nil {
		q.Set("dismissed", strconv.FormatBool(*in.Dismissed))
	}
	if in.Last != nil {
		q.Set("last", strconv.Itoa(int(*in.Last)))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	path += "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

func (c *Client) DismissAgentError(ctx context.Context, in types.DismissAgentError) (types.DismissedAgentError, error) {
	var out types.DismissedAgentError
	return out, c.do(ctx, http.MethodPost, "/v1/agent_errors/"+in.AgentErrorID+"/dismiss", in, &out)
}
