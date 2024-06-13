package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) CreateSAMLMapping(ctx context.Context, in types.CreateSAMLMapping) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/projects/"+url.PathEscape(in.ProjectID)+"/saml_mappings", in, &out)
}

func (c *Client) SAMLMappings(ctx context.Context, in types.ListSAMLMappings) (types.SAMLMappings, error) {
	var out types.SAMLMappings
	endpoint := "/v1/projects/" + url.PathEscape(in.ProjectID) + "/saml_mappings"
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

func (c *Client) UpdateSAMLMapping(ctx context.Context, in types.UpdateSAMLMapping) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/saml_mappings/"+url.PathEscape(in.ID), in, &out)
}

func (c *Client) DeleteSAMLMapping(ctx context.Context, id string) (types.Deleted, error) {
	var out types.Deleted
	return out, c.do(ctx, http.MethodDelete, "/v1/saml_mappings/"+url.PathEscape(id), nil, &out)
}
