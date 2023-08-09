package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateProcessingRuleTemplate(ctx context.Context, in types.CreateProcessingRuleTemplate) (types.Created, error) {
	var out types.Created
	path := "/v1/projects/" + url.PathEscape(in.ProjectID) + "/processing_rule_templates"
	return out, c.do(ctx, http.MethodPost, path, in, &out)
}

func (c *Client) ProcessingRuleTemplates(ctx context.Context, in types.ListProcessingRuleTemplates) (types.ProcessingRuleTemplates, error) {
	var out types.ProcessingRuleTemplates
	path := "/v1/projects/" + url.PathEscape(in.ProjectID) + "/processing_rule_templates"
	q := url.Values{}
	if in.Name != nil {
		q.Set("name", *in.Name)
	}
	if in.Last != nil {
		q.Set("last", strconv.Itoa(int(*in.Last)))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	path += "?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}

func (c *Client) UpdateProcessingRuleTemplate(ctx context.Context, in types.UpdateProcessingRuleTemplate) (types.Updated, error) {
	var out types.Updated
	path := "/v1/processing_rule_templates/" + url.PathEscape(in.TemplateID)
	return out, c.do(ctx, http.MethodPatch, path, in, &out)
}

func (c *Client) DeleteProcessingRuleTemplate(ctx context.Context, templateID string) (types.Deleted, error) {
	var out types.Deleted
	path := "/v1/processing_rule_templates/" + url.PathEscape(templateID)
	return out, c.do(ctx, http.MethodDelete, path, nil, &out)
}
