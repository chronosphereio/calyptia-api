package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateProcessingRule(ctx context.Context, in types.CreateProcessingRule) (types.CreatedProcessingRule, error) {
	var out types.CreatedProcessingRule
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(in.PipelineID)+"/processing_rules", in, &out)
}

func (c *Client) ProcessingRules(ctx context.Context, params types.ProcessingRulesParams) (types.ProcessingRules, error) {
	var out types.ProcessingRules
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}
	path := "/v1/pipelines/" + url.PathEscape(params.PipelineID) + "/processing_rules?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

func (c *Client) ProcessingRule(ctx context.Context, processingRuleID string) (types.ProcessingRule, error) {
	var out types.ProcessingRule
	return out, c.do(ctx, http.MethodGet, "/v1/processing_rules/"+url.PathEscape(processingRuleID), nil, &out)
}

func (c *Client) UpdateProcessingRule(ctx context.Context, in types.UpdateProcessingRule) (types.UpdatedProcessingRule, error) {
	var out types.UpdatedProcessingRule
	return out, c.do(ctx, http.MethodPatch, "/v1/processing_rules/"+url.PathEscape(in.ProcessingRuleID), in, &out)
}

func (c *Client) DeleteProcessingRule(ctx context.Context, processingRuleID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/processing_rules/"+url.PathEscape(processingRuleID), nil, nil)
}

func (c *Client) PreviewProcessingRule(ctx context.Context, in types.PreviewProcessingRule) ([]types.FluentBitLog, error) {
	var out []types.FluentBitLog
	return out, c.do(ctx, http.MethodPost, "/v1/preview_processing_rule", in, &out)
}
