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
	Lifespan   time.Duration `json:"lifespan" yaml:"lifespan"`
	CreatedAt  time.Time     `json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time     `json:"updatedAt" yaml:"updatedAt"`
}

// TraceSessionsParams request payload for querying trace sessions.
type TraceSessionsParams struct {
	Last   *uint64
	Before *string
}

// TraceRecord represents a single record in a trace session.
// Holds information about the plugin it represents
// and the trace record data.
type TraceRecord struct {
	ID        string          `json:"id" yaml:"id"`
	SessionID string          `json:"sessionID" yaml:"sessionID"`
	Plugin    string          `json:"plugin" yaml:"plugin"`
	Record    json.RawMessage `json:"record" yaml:"record"`
	CreatedAt time.Time       `json:"createdAt" yaml:"createdAt"`
}

// TraceRecordsParams request payload for querying trace records.
type TraceRecordsParams struct {
	Last   *uint64
	Before *string
}
