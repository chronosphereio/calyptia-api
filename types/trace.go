package types

import (
	"encoding/json"
	"time"
)

// TraceSession model.
type TraceSession struct {
	ID         string        `json:"id" yaml:"id"`
	PipelineID string        `json:"pipelineID" yaml:"pipelineID"`
	Plugins    []string      `json:"plugins" yaml:"plugins"`
	Lifespan   time.Duration `json:"lifespan" yaml:"lifespan"`
	CreatedAt  time.Time     `json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt" yaml:"updatedAt"`
}

// Active tells whether a session is still within its lifespan.
func (ts TraceSession) Active() bool {
	return ts.CreatedAt.Add(ts.Lifespan).After(time.Now())
}

// CreateTraceSession request payload for creating a new trace session.
type CreateTraceSession struct {
	Plugins  []string      `json:"tracePlugins"`
	Lifespan time.Duration `json:"traceLifespan"`
}

// TraceSessionsParams request payload for querying trace sessions.
type TraceSessionsParams struct {
	Last   *uint64
	Before *string
}

// TraceSessions paginated list.
type TraceSessions struct {
	Items     []TraceSession
	EndCursor *string
}

// UpdateTraceSession request payload for updating a trace session.
type UpdateTraceSession struct {
	Plugins  *[]string      `json:"plugins"`
	Lifespan *time.Duration `json:"lifespan"`
}

// TraceRecord model.
type TraceRecord struct {
	ID        string    `json:"id" yaml:"id"`
	SessionID string    `json:"sessionID" yaml:"sessionID"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`

	// fluent-bit data from here on.
	Kind           TraceRecordKind `json:"type" yaml:"type"`
	TraceID        string          `json:"traceID" yaml:"traceID"`
	StartTime      time.Time       `json:"start_time" yaml:"start_time"`
	EndTime        time.Time       `json:"end_time" yaml:"end_time"`
	InputInstance  string          `json:"input_instance" yaml:"input_instance"`
	FilterInstance string          `json:"filter_instance" yaml:"filter_instance"`
	OutputInstance string          `json:"outputInstance" yaml:"outputInstance"`
	ReturnCode     int             `json:"return_code" yaml:"return_code"`
	// Each record is a JSON object,
	// warranted to have a flb_time `timestamp` field.
	Records []json.RawMessage `json:"records" yaml:"records"`
}

// TraceRecordKind enum.
type TraceRecordKind string

const (
	TraceRecordKindInput     TraceRecordKind = "input"
	TraceRecordKindFilter    TraceRecordKind = "filter"
	TraceRecordKindPreOutput TraceRecordKind = "pre_output"
	TraceRecordKindOutput    TraceRecordKind = "output"
)

// CreateTraceRecord request payload for creating a new trace record.
type CreateTraceRecord struct {
	Kind           TraceRecordKind `json:"type" yaml:"type"`
	TraceID        string          `json:"traceID" yaml:"traceID"`
	StartTime      time.Time       `json:"start_time" yaml:"start_time"`
	EndTime        time.Time       `json:"end_time" yaml:"end_time"`
	InputInstance  string          `json:"input_instance" yaml:"input_instance"`
	FilterInstance string          `json:"filter_instance" yaml:"filter_instance"`
	OutputInstance string          `json:"outputInstance" yaml:"outputInstance"`
	ReturnCode     int             `json:"return_code" yaml:"return_code"`
	// Each record is a JSON object,
	// warranted to have a flb_time `timestamp` field.
	Records []json.RawMessage `json:"records" yaml:"records"`
}

// TraceRecordsParams request payload for querying trace records.
type TraceRecordsParams struct {
	Last   *uint64
	Before *string
}

// TraceRecords paginated list.
type TraceRecords struct {
	Items     []TraceRecord
	EndCursor *string
}
