package types

import (
	"encoding/json"
	"time"
)

// Pipeline model.
type Pipeline struct {
	ID              string           `json:"id" yaml:"id"`
	Name            string           `json:"name" yaml:"name"`
	Config          PipelineConfig   `json:"config" yaml:"config"`
	Status          PipelineStatus   `json:"status" yaml:"status"`
	ResourceProfile ResourceProfile  `json:"resourceProfile" yaml:"resourceProfile"`
	ReplicasCount   uint64           `json:"replicasCount" yaml:"replicasCount"`
	Metadata        *json.RawMessage `json:"metadata"`
	CreatedAt       time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

// CreatePipeline request payload for creating a new pipeline.
type CreatePipeline struct {
	Name                      string                 `json:"name"`
	ReplicasCount             uint64                 `json:"replicasCount"`
	RawConfig                 string                 `json:"rawConfig"`
	Secrets                   []CreatePipelineSecret `json:"secrets"`
	Files                     []CreatePipelineFile   `json:"files"`
	ResourceProfileName       string                 `json:"resourceProfile"`
	AutoCreatePortsFromConfig bool                   `json:"autoCreatePortsFromConfig"`
	Metadata                  *json.RawMessage       `json:"metadata"`
}

// CreatedPipeline response payload after creating a pipeline successfully.
type CreatedPipeline struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Config          PipelineConfig   `json:"config"`
	Secrets         []PipelineSecret `json:"secrets"`
	Files           []PipelineFile   `json:"files"`
	Status          PipelineStatus   `json:"status"`
	ResourceProfile ResourceProfile  `json:"resourceProfile"`
	ReplicasCount   uint64           `json:"replicasCount"`
	CreatedAt       time.Time        `json:"createdAt"`
}

// UpdatePipeline request payload for updating a pipeline.
type UpdatePipeline struct {
	Name                      *string                `json:"name"`
	ReplicasCount             *uint64                `json:"replicasCount"`
	RawConfig                 *string                `json:"rawConfig"`
	Secrets                   []UpdatePipelineSecret `json:"secrets"`
	Files                     []UpdatePipelineFile   `json:"files"`
	Status                    *PipelineStatusKind    `json:"status"`
	ResourceProfile           *string                `json:"resourceProfile"`
	AutoCreatePortsFromConfig *bool                  `json:"autoCreatePortsFromConfig"`
	Metadata                  *json.RawMessage       `json:"metadata"`
}

// PipelinesParams request payload for querying pipelines.
type PipelinesParams struct {
	Last *uint64
	Name *string
}

// UpdatedPipeline response payload after updating a pipeline successfully.
type UpdatedPipeline struct {
	AddedPorts []PipelinePort `json:"addedPorts"`
}
