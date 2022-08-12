package types

import (
	"encoding/json"
	"time"
)

// TraceSession model.
type TraceSession struct {
	ID         string    `json:"id" yaml:"id"`
	PipelineID string    `json:"pipelineID" yaml:"pipelineID"`
	Plugins    []string  `json:"plugins" yaml:"plugins"`
	EndTime    time.Time `json:"endTime" yaml:"endTime"`
	CreatedAt  time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// Active tells whether a session is still within its lifespan.
func (ts TraceSession) Active() bool {
	now := time.Now()
	return ts.EndTime.Equal(now) || ts.EndTime.After(now)
}

// CreateTraceSession request payload for creating a new trace session.
type CreateTraceSession struct {
	Plugins []string  `json:"plugins"`
	EndTime time.Time `json:"endTime"`
}

// TraceSessionsParams request payload for querying trace sessions.
type TraceSessionsParams struct {
	Last   *uint64
	Before *string
}

// TraceSessions paginated list.
type TraceSessions struct {
	Items     []TraceSession `json:"items"`
	EndCursor *string        `json:"endCursor"`
}

// UpdateTraceSession request payload for updating a trace session.
type UpdateTraceSession struct {
	Plugins *[]string  `json:"plugins"`
	EndTime *time.Time `json:"endTime"`
}

// TerminatedTraceSession response payload after terminating a trace session.
type TerminatedTraceSession struct {
	ID        string    `json:"id" yaml:"id"`
	EndTime   time.Time `json:"endTime" yaml:"endTime"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// TraceRecord model.
type TraceRecord struct {
	ID        string    `json:"id" yaml:"id"`
	SessionID string    `json:"sessionID" yaml:"sessionID"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`

	// fluent-bit data from here on.
	Kind           TraceRecordKind `json:"type" yaml:"type"`
	TraceID        string          `json:"trace_id" yaml:"trace_id"`
	StartTime      time.Time       `json:"start_time" yaml:"start_time"`
	EndTime        time.Time       `json:"end_time" yaml:"end_time"`
	PluginInstance string          `json:"plugin_instance" yaml:"plugin_instance"`
	PluginAlias    string          `json:"plugin_alias" yaml:"plugin_alias"`
	ReturnCode     int             `json:"return_code" yaml:"return_code"`
	// Records array, each record is a JSON object,
	// warranted to have a flb_time `timestamp` field.
	Records json.RawMessage `json:"records" yaml:"records"`
}

// TraceRecordKind enum.
type TraceRecordKind uint

const (
	TraceRecordKindInput TraceRecordKind = iota + 1
	TraceRecordKindFilter
	TraceRecordKindPreOutput
	TraceRecordKindOutput
)

// CreateTraceRecord request payload for creating a new trace record.
type CreateTraceRecord struct {
	Kind    TraceRecordKind `json:"type"`
	TraceID string          `json:"trace_id"`
	// StartTime in unix seconds.
	StartTime int64 `json:"start_time"`
	// EndTime in unix seconds.
	EndTime        int64  `json:"end_time"`
	PluginInstance string `json:"plugin_instance"`
	PluginAlias    string `json:"plugin_alias"`
	ReturnCode     int    `json:"return_code"`
	// Records array, each record is a JSON object,
	// warranted to have a flb_time `timestamp` field.
	Records json.RawMessage `json:"records"`
}

// TraceRecordsParams request payload for querying trace records.
type TraceRecordsParams struct {
	Last   *uint64
	Before *string
}

// TraceRecords paginated list.
type TraceRecords struct {
	Items     []TraceRecord `json:"items"`
	EndCursor *string       `json:"endCursor"`
}
