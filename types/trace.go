package types

import (
	"encoding/json"
	"time"
)

// TraceSession gets created alongside a pipeline
// when tracing is enabled.
// This will run for the configured lifespan
// and process records during that time.
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

// TraceRecord represents a single record in a trace session.
// Holds information about the trace record data.
type TraceRecord struct {
	ID        string          `json:"id" yaml:"id"`
	SessionID string          `json:"sessionID" yaml:"sessionID"`
	Record    json.RawMessage `json:"record" yaml:"record"`
	CreatedAt time.Time       `json:"createdAt" yaml:"createdAt"`
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
