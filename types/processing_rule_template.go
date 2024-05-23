package types

import (
	"encoding/json"
	"time"
)

type ProcessingRuleTemplate struct {
	ID              string          `json:"id" yaml:"id" db:"id"`
	ProjectID       string          `json:"projectID" yaml:"projectID" db:"project_id"`
	Name            string          `json:"name" yaml:"name" db:"name"`
	PipelineVersion *string         `json:"pipelineVersion" yaml:"pipelineVersion" db:"pipeline_version"`
	Definition      json.RawMessage `json:"definition" yaml:"definition" db:"definition"`
	CreatedAt       time.Time       `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt       time.Time       `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreateProcessingRuleTemplate struct {
	ProjectID       string          `json:"-"`
	Name            string          `json:"name"`
	PipelineVersion *string         `json:"pipelineVersion"`
	Definition      json.RawMessage `json:"definition"`
}

type ListProcessingRuleTemplates struct {
	ProjectID string  `json:"-"`
	Last      *uint   `json:"-"`
	Before    *string `json:"-"`
	Name      *string `json:"-"`
}

type ProcessingRuleTemplates struct {
	Items     []ProcessingRuleTemplate `json:"items" yaml:"items"`
	EndCursor *string                  `json:"endCursor" yaml:"endCursor"`
	Count     uint                     `json:"count" yaml:"count"`
}

type UpdateProcessingRuleTemplate struct {
	ID              string           `json:"-"`
	Name            *string          `json:"name"`
	PipelineVersion *string          `json:"pipelineVersion"`
	Definition      *json.RawMessage `json:"definition"`
}
