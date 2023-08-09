package client

import (
	"context"
	"net/http"
	"net/url"
	"strconv"

	"github.com/calyptia/api/types"
)

// CreateTraceSession within a pipeline.
// A trace session can receive records from fluent-bit's tracing sidecar
// while this is enabled (see its lifespan).
// Only one trace session can be active at a time
// you can either terminate it and create a new one, or update the existing one
// and extend its lifespan.
func (c *Client) CreateTraceSession(ctx context.Context, pipelineID string, in types.CreateTraceSession) (types.Created, error) {
	var out types.Created
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/trace_sessions", in, &out)
}

// TraceSessions from a pipeline.
// With backward pagination, the list is sorted by createdAt in descending order.
func (c *Client) TraceSessions(ctx context.Context, pipelineID string, params types.TraceSessionsParams) (types.TraceSessions, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.TraceSessions
	path := "/v1/pipelines/" + url.PathEscape(pipelineID) + "/trace_sessions?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}

// TraceSession by ID. This can be still active or not.
func (c *Client) TraceSession(ctx context.Context, sessionID string) (types.TraceSession, error) {
	var out types.TraceSession
	return out, c.do(ctx, http.MethodGet, "/v1/trace_sessions/"+url.PathEscape(sessionID), nil, &out)
}

// ActiveTraceSession from a pipeline if any.
func (c *Client) ActiveTraceSession(ctx context.Context, pipelineID string) (types.TraceSession, error) {
	var out types.TraceSession
	return out, c.do(ctx, http.MethodGet, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/trace_session", nil, &out)
}

// UpdateTraceSession list of plugins to trace and/or lifespan.
// The session must still be active.
func (c *Client) UpdateTraceSession(ctx context.Context, sessionID string, in types.UpdateTraceSession) (types.Updated, error) {
	var out types.Updated
	return out, c.do(ctx, http.MethodPatch, "/v1/trace_sessions/"+url.PathEscape(sessionID), in, &out)
}

// TerminateActiveTraceSession terminates the current active trace session
// on the pipeline.
func (c *Client) TerminateActiveTraceSession(ctx context.Context, pipelineID string) (types.TerminatedTraceSession, error) {
	var out types.TerminatedTraceSession
	return out, c.do(ctx, http.MethodDelete, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/trace_session", nil, &out)
}

// CreateTraceRecord on the current active trace session on the pipeline.
// This endpoint is meant to be used by fluent-bit's tracing sidecar.
func (c *Client) CreateTraceRecord(ctx context.Context, pipelineID string, in types.CreateTraceRecord) (types.CreatedTraceRecord, error) {
	var out types.CreatedTraceRecord
	return out, c.do(ctx, http.MethodPost, "/v1/pipelines/"+url.PathEscape(pipelineID)+"/trace_session/records", in, &out)
}

// TraceRecords from a trace session.
// With backward pagination, the list is sorted by createdAt in descending order.
func (c *Client) TraceRecords(ctx context.Context, sessionID string, params types.TraceRecordsParams) (types.TraceRecords, error) {
	q := url.Values{}
	if params.Last != nil {
		q.Set("last", strconv.FormatUint(uint64(*params.Last), uintBase))
	}
	if params.Before != nil {
		q.Set("before", *params.Before)
	}

	var out types.TraceRecords
	path := "/v1/trace_sessions/" + url.PathEscape(sessionID) + "/records?" + q.Encode()
	return out, c.do(ctx, http.MethodGet, path, nil, &out.Items, withCursor(&out.EndCursor))
}
