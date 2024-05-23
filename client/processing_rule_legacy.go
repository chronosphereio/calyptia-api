package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
	"github.com/calyptia/api/types/legacy"
)

func (c *Client) CreateProcessingRuleV1(ctx context.Context, in legacy.CreateProcessingRule) (legacy.CreatedProcessingRule, error) {
	var out legacy.CreatedProcessingRule
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(in.PipelineID)+"/processing_rules", in, &out)
}

func (c *Client) ProcessingRulesV1(ctx context.Context, params legacy.ProcessingRulesParams) (legacy.ProcessingRules, error) {
	var out legacy.ProcessingRules
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

func (c *Client) ProcessingRuleV1(ctx context.Context, processingRuleID string) (legacy.ProcessingRule, error) {
	var out legacy.ProcessingRule
	return out, c.do(ctx, http.MethodGet, "/v1/processing_rules/"+url.PathEscape(processingRuleID), nil, &out)
}

func (c *Client) UpdateProcessingRuleV1(ctx context.Context, in legacy.UpdateProcessingRule) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/processing_rules/"+url.PathEscape(in.ProcessingRuleID), in, &out)
}

func (c *Client) DeleteProcessingRuleV1(ctx context.Context, processingRuleID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/processing_rules/"+url.PathEscape(processingRuleID), nil, nil)
}

func (c *Client) PreviewProcessingRuleV1(ctx context.Context, in legacy.PreviewProcessingRule) ([]types.FluentBitLog, error) {
	var out []types.FluentBitLog
	return out, c.do(ctx, http.MethodPost, "/v1/preview_processing_rule", in, &out)
}
