package types

import "time"

type ConfigFormat string

const (
	ConfigFormatINI  ConfigFormat = "ini"
	ConfigFormatJSON ConfigFormat = "json"
	ConfigFormatYAML ConfigFormat = "yaml"
)

// PipelineConfig model.
type PipelineConfig struct {
	ID           string       `json:"id" yaml:"id"`
	RawConfig    string       `json:"rawConfig" yaml:"rawConfig"`
	ConfigFormat ConfigFormat `json:"configFormat"`
	CreatedAt    time.Time    `json:"createdAt" yaml:"createdAt"`
}

// PipelineConfigHistory paginated list.
type PipelineConfigHistory struct {
	Items     []PipelineConfig
	EndCursor *string
}

// PipelineConfigHistoryParams request payload for querying the pipeline config history.
type PipelineConfigHistoryParams struct {
	Last         *uint
	Before       *string
	ConfigFormat *ConfigFormat
}
