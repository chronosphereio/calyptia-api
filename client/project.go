package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateProject creates a new project.
// A project is the base unit of work at Calyptia Cloud.
// You can register agents here, create aggregators in which you can deploy an entire set of pipelines, and monitor them.
// You can even invite other people to the project and have a team.
func (c *Client) CreateProject(ctx context.Context, payload types.CreateProject) (types.CreatedProject, error) {
	var out types.CreatedProject
	return out, c.do(ctx, http.MethodPost, "/v1/projects", payload, &out)
}

// Projects you are a member of in descending order.
func (c *Client) Projects(ctx context.Context, params types.ProjectsParams) (types.Projects, error) {
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

	var out types.Projects
	return out, c.do(ctx, http.MethodGet, "/v1/projects?"+q.Encode(), nil, &out.Items, withCursor(&out.EndCursor))
}

// Project by ID.
func (c *Client) Project(ctx context.Context, projectID string) (types.Project, error) {
	var out types.Project
	return out, c.do(ctx, http.MethodGet, "/v1/projects/"+url.PathEscape(projectID), nil, &out)
}

// UpdateProject by its ID.
func (c *Client) UpdateProject(ctx context.Context, projectID string, opts types.UpdateProject) error {
	return c.do(ctx, http.MethodPatch, "/v1/projects/"+url.PathEscape(projectID), opts, nil)
}
