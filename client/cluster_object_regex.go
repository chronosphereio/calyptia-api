package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateClusterObjectRegex(ctx context.Context, in types.CreateClusterObjectRegex) (types.CreatedClusterObjectRegex, error) {
	var out types.CreatedClusterObjectRegex
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(in.PipelineID)+"/cluster_object_regexes", in, &out)
}

func (c *Client) ClusterObjectRegexes(ctx context.Context, in types.ListClusterObjectRegexes) (types.ClusterObjectRegexes, error) {
	var out types.ClusterObjectRegexes
	q := url.Values{}
	if in.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*in.Last), uintBase))
	}
	if in.Before != nil {
		q.Set("before", *in.Before)
	}
	path := "/v1/pipelines/" + url.PathEscape(in.PipelineID) + "/cluster_object_regexes"
	if len(q) > 0 {
		path += "?" + q.Encode()
	}
	return out, c.do(ctx, http.MethodGet, path, in, &out)
}

func (c *Client) ClusterObjectRegex(ctx context.Context, regexID string) (types.ClusterObjectRegex, error) {
	var out types.ClusterObjectRegex
	return out, c.do(ctx, http.MethodGet, "/v1/cluster_object_regexes/"+url.PathEscape(regexID), nil, &out)
}

func (c *Client) UpdateClusterObjectRegex(ctx context.Context, in types.UpdateClusterObjectRegex) (types.UpdatedClusterObjectRegex, error) {
	var out types.UpdatedClusterObjectRegex
	return out, c.do(ctx, http.MethodPatch, "/v1/cluster_object_regexes/"+url.PathEscape(in.RegexID), in, &out)
}

func (c *Client) DeleteClusterObjectRegex(ctx context.Context, regexID string) error {
	return c.do(ctx, http.MethodDelete, "/v1/cluster_object_regexes/"+url.PathEscape(regexID), nil, nil)
}
