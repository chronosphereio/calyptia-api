package types

import "time"

// PipelineStatus model.
type PipelineStatus struct {
	ID        string             `json:"id" yaml:"id"`
	Config    PipelineConfig     `json:"config" yaml:"config"`
	Status    PipelineStatusKind `json:"status" yaml:"status"`
	CreatedAt time.Time          `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" yaml:"updatedAt"`
}

type PipelineStatusKind string

const (
	PipelineStatusNew      PipelineStatusKind = "NEW"
	PipelineStatusFailed   PipelineStatusKind = "FAILED"
	PipelineStatusStarting PipelineStatusKind = "STARTING"
	PipelineStatusStarted  PipelineStatusKind = "STARTED"
	PipelineStatusDeleted  PipelineStatusKind = "DELETED"
)

// PipelineStatusHistoryParams request payload for querying the pipeline status history.
type PipelineStatusHistoryParams struct {
	Last *uint64
}
