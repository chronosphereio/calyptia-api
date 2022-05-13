package types

import "time"

// PipelineConfig model.
type PipelineConfig struct {
	ID           string    `json:"id" yaml:"id"`
	ConfigFormat *string   `json:"configFormat"`
	RawConfig    string    `json:"rawConfig" yaml:"rawConfig"`
	CreatedAt    time.Time `json:"createdAt" yaml:"createdAt"`
}

// PipelineConfigHistory paginated list.
type PipelineConfigHistory struct {
	Items     []PipelineConfig
	EndCursor *string
}

// PipelineConfigHistoryParams request payload for querying the pipeline config history.
type PipelineConfigHistoryParams struct {
	Last   *uint64
	Before *string
}
