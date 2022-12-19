package types

import (
	"encoding/json"
	"errors"
	"time"
)

// TraceSession model.
type TraceSession struct {
	ID         string `json:"id" yaml:"id"`
	PipelineID string `json:"pipelineID" yaml:"pipelineID"`
	// Plugins list of fluent-bit plugin IDs or aliases. Example ["forward.0", "myplugin"]
	Plugins   []string  `json:"plugins" yaml:"plugins"`
	Lifespan  Duration  `json:"lifespan" yaml:"lifespan"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// Active tells whether a session is still within its lifespan.
func (ts TraceSession) Active() bool {
	now := time.Now().Add(time.Microsecond)
	return ts.EndTime().Equal(now) || ts.EndTime().After(now)
}

// EndTime returns the time when the session will stop being active.
func (ts TraceSession) EndTime() time.Time {
	return ts.CreatedAt.Add(time.Duration(ts.Lifespan))
}

// CreateTraceSession request payload for creating a new trace session.
type CreateTraceSession struct {
	Plugins  []string `json:"plugins"`
	Lifespan Duration `json:"lifespan"`
}

// CreatedTraceSession response payload after creating a trace session successfully.
type CreatedTraceSession struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// TraceSessionsParams request payload for querying trace sessions.
type TraceSessionsParams struct {
	Last   *uint
	Before *string
}

// TraceSessions paginated list.
type TraceSessions struct {
	Items     []TraceSession `json:"items"`
	EndCursor *string        `json:"endCursor"`
}

// UpdateTraceSession request payload for updating a trace session.
type UpdateTraceSession struct {
	Plugins  *[]string `json:"plugins"`
	Lifespan *Duration `json:"lifespan"`
}

// UpdatedTraceSession response payload after updating a trace session successfully.
type UpdatedTraceSession struct {
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// TerminatedTraceSession response payload after terminating the active trace session successfully.
type TerminatedTraceSession struct {
	ID        string    `json:"id" yaml:"id"`
	Lifespan  Duration  `json:"lifespan" yaml:"lifespan"`
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

// String implements fmt.Stringer.
func (k TraceRecordKind) String() string {
	switch k {
	case TraceRecordKindInput:
		return "input"
	case TraceRecordKindFilter:
		return "filter"
	case TraceRecordKindPreOutput:
		return "pre_output"
	case TraceRecordKindOutput:
		return "output"
	default:
		return ""
	}
}

// GoString implements fmt.GoStringer.
func (k TraceRecordKind) GoString() string { return k.String() }

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

// CreatedTraceRecord response payload after creating an session record successfully.
type CreatedTraceRecord struct {
	ID        string    `json:"id" yaml:"id"`
	SessionID string    `json:"sessionID" yaml:"sessionID"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// TraceRecordsParams request payload for querying trace records.
type TraceRecordsParams struct {
	Last   *uint
	Before *string
}

// TraceRecords paginated list.
type TraceRecords struct {
	Items     []TraceRecord `json:"items"`
	EndCursor *string       `json:"endCursor"`
}

// Duration is a time.Duration wrapper that adds support for encoding/json.
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case int64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

func (d Duration) AsDuration() time.Duration {
	return time.Duration(d)
}
