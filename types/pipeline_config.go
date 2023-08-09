package types

import (
	"time"

	"github.com/calyptia/api/types/errs"
	fluentbitconfig "github.com/calyptia/go-fluentbit-config/v2"
)

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

func (c *PipelineConfig) ApplyFormat(format ConfigFormat) error {
	if c.ConfigFormat == "" {
		c.ConfigFormat = ConfigFormatINI
	}

	if c.ConfigFormat == format {
		return nil
	}

	parsed, err := fluentbitconfig.ParseAs(c.RawConfig, fluentbitconfig.Format(c.ConfigFormat))
	if err != nil {
		return errs.NewDetailedError(errs.InvalidPipelineConfig, err.Error())
	}

	raw, err := parsed.DumpAs(fluentbitconfig.Format(format))
	if err != nil {
		return errs.NewDetailedError(errs.InvalidPipelineConfig, err.Error())
	}

	c.RawConfig = raw
	c.ConfigFormat = format

	return nil
}

type CreatePipelineConfig struct {
	RawConfig    string       `json:"rawConfig"`
	ConfigFormat ConfigFormat `json:"configFormat"`
}

// PipelineConfigHistoryParams request payload for querying the pipeline config history.
type PipelineConfigHistoryParams struct {
	Last         *uint
	Before       *string
	ConfigFormat *ConfigFormat
}

// PipelineConfigHistory paginated list.
type PipelineConfigHistory struct {
	Items     []PipelineConfig
	EndCursor *string
}
