package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateProject(ctx context.Context, payload types.CreateProject) (types.CreatedProject, error) {
	var out types.CreatedProject
	return out, c.do(ctx, http.MethodPost, "/v1/projects", payload, &out)
}

func (c *Client) Projects(ctx context.Context, params types.ProjectsParams) ([]types.Project, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(*params.Last, 10))
	}
	if params.Name != nil {
		q.Set("name", *params.Name)
	}

	var out []types.Project
	return out, c.do(ctx, http.MethodGet, "/v1/projects?"+q.Encode(), nil, &out)
}

func (c *Client) Project(ctx context.Context, projectID string) (types.Project, error) {
	var out types.Project
	return out, c.do(ctx, http.MethodGet, "/v1/projects/"+url.PathEscape(projectID), nil, &out)
}

func (c *Client) UpdateProject(ctx context.Context, projectID string, opts types.UpdateProject) error {
	return c.do(ctx, http.MethodPatch, "/v1/projects/"+url.PathEscape(projectID), opts, nil)
}
