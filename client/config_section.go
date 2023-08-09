package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateConfigSection(ctx context.Context, projectID string, in types.CreateConfigSection) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(projectID)+"/config_sections", in, &out)
}

func (c *Client) ConfigSections(ctx context.Context, projectID string, params types.ConfigSectionsParams) (types.ConfigSections, error) {
	var out types.ConfigSections

	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	if params.IncludeProcessingRules {
		q.Set("include_processing_rules", "true")
	}

	err := c.do(ctx, http.MethodGet, "/v1/projects/"+url.PathEscape(projectID)+"/config_sections?"+q.Encode(), nil, &out.Items, withCursor(&out.EndCursor))
	if err != nil {
		return out, err
	}

	for _, sc := range out.Items {
		sc.Properties.FixFloats()
	}

	return out, nil
}

func (c *Client) ConfigSection(ctx context.Context, configSectionID string) (types.ConfigSection, error) {
	var out types.ConfigSection
	err := c.do(ctx, http.MethodGet, "/v1/config_sections/"+url.PathEscape(configSectionID), nil, &out)
	if err != nil {
		return out, err
	}

	out.Properties.FixFloats()

	return out, nil
}

func (c *Client) UpdateConfigSection(ctx context.Context, configSectionID string, in types.UpdateConfigSection) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/config_sections/"+url.PathEscape(configSectionID), in, &out)
}

func (c *Client) DeleteConfigSection(ctx context.Context, configSectionID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/config_sections/"+url.PathEscape(configSectionID), nil, nil)
}

func (c *Client) UpdateConfigSectionSet(ctx context.Context, pipelineID string, configSectionIDs ...string) error {
	return c.do(ctx, http.MethodPatch, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/config_section_set", configSectionIDs, nil)
}
