package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// RegisterAgent within a project.
// The project in which the agent is registered is parsed from the authorization token.
// Users are not allowed to register agents.
func (c *Client) RegisterAgent(ctx context.Context, payload types.RegisterAgent) (types.RegisteredAgent, error) {
	var out types.RegisteredAgent
	return out, c.do(ctx, http.MethodPost, "/v1/agents", payload, &out)
}

// Agents from the given project in descending order.
func (c *Client) Agents(ctx context.Context, projectID string, params types.AgentsParams) (types.Agents, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}
	if params.TagsQuery != nil {
		q.Set("tags_query", *params.TagsQuery)
	}
	if params.EnvironmentID != nil {
		q.Set("environment_id", *params.EnvironmentID)
	}
	if params.FleetID != nil {
		q.Set("fleet_id", *params.FleetID)
	}

	var out types.Agents
	path := "/v1/projects/" + url.PathEscape(projectID) + "/agents?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// Agent by ID.
func (c *Client) Agent(ctx context.Context, agentID string) (types.Agent, error) {
	var out types.Agent
	return out, c.do(ctx, http.MethodGet, "/v1/agents/"+url.PathEscape(agentID), nil, &out)
}

// UpdateAgent by its ID.
// Users are allowed to only update a restricted set of fields (name);
// while agents are allowed to update the whole configuration.
func (c *Client) UpdateAgent(ctx context.Context, agentID string, payload types.UpdateAgent) error {
	return c.do(ctx, http.MethodPatch, "/v1/agents/"+url.PathEscape(agentID), payload, nil)
}

func (c *Client) AddAgentMetrics(ctx context.Context, agentID string, payload []byte) error {
	return c.do(ctx, http.MethodPost, "/v1/agents/"+url.PathEscape(agentID)+"/metrics", payload, nil)
}

// DeleteAgent by its ID.
func (c *Client) DeleteAgent(ctx context.Context, agentID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/agents/"+url.PathEscape(agentID), nil, nil)
}

// DeleteAgents from a project passing a list of the IDs to be deleted.
func (c *Client) DeleteAgents(ctx context.Context, projectID string, agentIDs ...string) error {
	var errs error
	for _, chunk := range makeChunks(agentIDs, chunksSizeByID) {
		q := url.Values{}
		for _, id := range chunk {
			q.Add("agent_id", id)
		}
		err := c.do(ctx, http.MethodDelete, "/v1/projects/"+url.PathEscape(projectID)+"/agents?"+q.Encode(), nil, nil)
		if err != nil {
			errs = errors.Join(errs, err)
		}
	}
	return errs
}
