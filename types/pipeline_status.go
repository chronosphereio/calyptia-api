package types

import "time"

// PipelineStatus model.
type PipelineStatus struct {
	ID        string             `json:"id" yaml:"id"`
	Config    PipelineConfig     `json:"config" yaml:"config"`
	Status    PipelineStatusKind `json:"status" yaml:"status"`
	CreatedAt time.Time          `json:"createdAt" yaml:"createdAt"`
}

// PipelineStatusKind enum.
type PipelineStatusKind string

const (
	// PipelineStatusNew is the default status right after a pipeline is created.
	PipelineStatusNew PipelineStatusKind = "NEW"
	// PipelineStatusFailed is the status when a pipeline fails.
	PipelineStatusFailed PipelineStatusKind = "FAILED"
	// PipelineStatusStarting is the status of a starting pipeline.
	PipelineStatusStarting PipelineStatusKind = "STARTING"
	// PipelineStatusStarted is the status of a started pipeline.
	PipelineStatusStarted PipelineStatusKind = "STARTED"
	// PipelineStatusDeleted is the status of a deleted pipeline.
	PipelineStatusDeleted PipelineStatusKind = "DELETED"
)

// PipelineStatusHistoryParams request payload for querying the pipeline status history.
type PipelineStatusHistoryParams struct {
	Last *uint64
}
