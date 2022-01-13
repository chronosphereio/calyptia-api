package types

import "time"

// PipelineStatus model.
type PipelineStatus struct {
	ID        string         `json:"id" yaml:"id"`
	Config    PipelineConfig `json:"config" yaml:"config"`
	Status    string         `json:"status" yaml:"status"`
	CreatedAt time.Time      `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" yaml:"updatedAt"`
}

// PipelineStatusHistoryParams request payload for querying the pipeline status history.
type PipelineStatusHistoryParams struct {
	Last *uint64
}
