package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) RegisterAgent(ctx context.Context, payload types.RegisterAgent) (types.RegisteredAgent, error) {
	var out types.RegisteredAgent
	return out, c.do(ctx, http.MethodPost, "/v1/agents", payload, &out)
}

func (c *Client) Agents(ctx context.Context, projectID string, params types.AgentsParams) ([]types.Agent, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out []types.Agent
	path := "/v1/projects/" + url.PathEscape(projectID) + "/agents?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) Agent(ctx context.Context, agentID string) (types.Agent, error) {
	var out types.Agent
	return out, c.do(ctx, http.MethodGet, "/v1/agents/"+url.PathEscape(agentID), nil, &out)
}

func (c *Client) UpdateAgent(ctx context.Context, agentID string, payload types.UpdateAgent) error {
	return c.do(ctx, http.MethodPatch, "/v1/agents/"+url.PathEscape(agentID), payload, nil)
}

func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/agents/"+url.PathEscape(agentID), nil, nil)
}
