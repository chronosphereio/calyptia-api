package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

func (c *Client) AuditLogs(ctx context.Context, params types.ListAuditLogs) (types.AuditLogs, error) {
	q := url.Values{}
	if params.First != nil {
		q.Set("first", strconv.FormatUint(uint64(*params.First), uintBase))
	}
	if params.After != nil {
		q.Set("after", *params.After)
	}

	if params.URL != "" {
		q.Set("url", params.URL)
	}
	if params.Action != "" {
		q.Set("action", params.URL)
	}

	for key, value := range params.Components {
		q.Set("component["+key+"]", value)
	}
	for key, value := range params.Identity {
		q.Set("identity["+key+"]", value)
	}

	var out types.AuditLogs

	path := "/v1/projects/" +
		url.PathEscape(params.ProjectID) +
		"/audit_logs?" + q.Encode()

	return out, c.do(ctx, http.MethodGet, path, nil, &out)
}
