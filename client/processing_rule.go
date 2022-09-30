package client

import (
	"context"
	"net/http"
	"net/url"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateProcessingRule(ctx context.Context, in types.CreateProcessingRule) (types.CreatedProcessingRule, error) {
	var out types.CreatedProcessingRule
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(in.PipelineID)+"/processing_rules", in, &out)
}

func (c *Client) ProcessingRuleFromPipeline(ctx context.Context, pipelineID string) (types.ProcessingRule, error) {
	var out types.ProcessingRule
	return out, c.do(ctx, http.MethodGet, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/processing_rule", nil, &out)
}

func (c *Client) UpdateProcessingRule(ctx context.Context, in types.UpdateProcessingRule) (types.UpdatedProcessingRule, error) {
	var out types.UpdatedProcessingRule
	return out, c.do(ctx, http.MethodPatch, "/v1/processing_rules/"+url.PathEscape(in.ProcessingRuleID), in, &out)
}
