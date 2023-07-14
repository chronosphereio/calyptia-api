package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateProcessingRuleTemplate(ctx context.Context, in types.CreateProcessingRuleTemplate) (types.CreatedProcessingRuleTemplate, error) {
	var out types.CreatedProcessingRuleTemplate
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

func (c *Client) UpdateProcessingRuleTemplate(ctx context.Context, in types.UpdateProcessingRuleTemplate) (types.UpdatedProcessingRuleTemplate, error) {
	var out types.UpdatedProcessingRuleTemplate
	path := "/v1/processing_rule_templates/" + url.PathEscape(in.TemplateID)
	return out, c.do(ctx, http.MethodPatch, path, in, &out)
}

func (c *Client) DeleteProcessingRuleTemplate(ctx context.Context, templateID string) (types.DeletedProcessingRuleTemplate, error) {
	var out types.DeletedProcessingRuleTemplate
	path := "/v1/processing_rule_templates/" + url.PathEscape(templateID)
	return out, c.do(ctx, http.MethodDelete, path, nil, &out)
}
