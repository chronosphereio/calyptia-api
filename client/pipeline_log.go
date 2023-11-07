package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreatePipelineLog(ctx context.Context, in types.CreatePipelineLog) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(in.PipelineID)+"/logs", in, &out)
}

func (c *Client) PipelineLogs(ctx context.Context, in types.ListPipelineLogs) (types.PipelineLogs, error) {
	var out types.PipelineLogs
	endpoint := "/v1/pipelines/" + url.PathEscape(in.PipelineID) + "/logs"
	q := url.Values{}
	if in.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*in.Last), 10))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	if len(q) != 0 {
		endpoint += "?" + q.Encode()
	}
	return out, c.do(ctx, http.MethodGet, endpoint, nil, &out)
}

func (c *Client) UpdatePipelineLog(ctx context.Context, in types.UpdatePipelineLog) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/pipeline_logs/"+url.PathEscape(in.ID), in, &out)
}

func (c *Client) DeletePipelineLog(ctx context.Context, id string) (types.Deleted, error) {
	var out types.Deleted
	return out, c.do(ctx, http.MethodDelete, "/v1/pipeline_logs/"+url.PathEscape(id), nil, &out)
}
