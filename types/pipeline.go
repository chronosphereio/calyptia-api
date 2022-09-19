package types

import (
	"encoding/json"
	"time"
)

type PipelineKind string

const (
	PipelineKindDaemonSet  PipelineKind = "daemonSet"
	PipelineKindDeployment PipelineKind = "deployment"
)

// Pipeline model.
type Pipeline struct {
	ID              string           `json:"id" yaml:"id"`
	Name            string           `json:"name" yaml:"name"`
	Kind            PipelineKind     `json:"kind" yaml:"kind"`
	Config          PipelineConfig   `json:"config" yaml:"config"`
	Status          PipelineStatus   `json:"status" yaml:"status"`
	ResourceProfile ResourceProfile  `json:"resourceProfile" yaml:"resourceProfile"`
	TracingEnabled  bool             `json:"tracingEnabled" yaml:"tracingEnabled"`
	ReplicasCount   uint             `json:"replicasCount" yaml:"replicasCount"`
	Tags            []string         `json:"tags" yaml:"tags"`
	Metadata        *json.RawMessage `json:"metadata" yaml:"metadata"`
	CreatedAt       time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

// Pipelines paginated list.
type Pipelines struct {
	Items     []Pipeline
	EndCursor *string
}

// CreatePipeline request payload for creating a new pipeline.
type CreatePipeline struct {
	Name                      string                 `json:"name"`
	Kind                      PipelineKind           `json:"kind"`
	ReplicasCount             uint                   `json:"replicasCount"`
	RawConfig                 string                 `json:"rawConfig"`
	ConfigFormat              ConfigFormat           `json:"configFormat"`
	Secrets                   []CreatePipelineSecret `json:"secrets"`
	Files                     []CreatePipelineFile   `json:"files"`
	ResourceProfileName       string                 `json:"resourceProfile"`
	AutoCreatePortsFromConfig bool                   `json:"autoCreatePortsFromConfig"`
	Metadata                  *json.RawMessage       `json:"metadata"`
	Tags                      []string               `json:"tags"`
}

// CreatedPipeline response payload after creating a pipeline successfully.
type CreatedPipeline struct {
	ID              string           `json:"id"`
	Name            string           `json:"name"`
	Kind            PipelineKind     `json:"kind"`
	Config          PipelineConfig   `json:"config"`
	Secrets         []PipelineSecret `json:"secrets"`
	Files           []PipelineFile   `json:"files"`
	Status          PipelineStatus   `json:"status"`
	ResourceProfile ResourceProfile  `json:"resourceProfile"`
	ReplicasCount   uint             `json:"replicasCount"`
	CreatedAt       time.Time        `json:"createdAt"`
}

// UpdatePipeline request payload for updating a pipeline.
type UpdatePipeline struct {
	Name                      *string                `json:"name"`
	Kind                      *PipelineKind          `json:"kind"`
	ReplicasCount             *uint                  `json:"replicasCount"`
	RawConfig                 *string                `json:"rawConfig"`
	ConfigFormat              *ConfigFormat          `json:"configFormat"`
	Secrets                   []UpdatePipelineSecret `json:"secrets"`
	Files                     []UpdatePipelineFile   `json:"files"`
	Status                    *PipelineStatusKind    `json:"status"`
	Events                    []PipelineEvent        `json:"events"`
	ResourceProfile           *string                `json:"resourceProfile"`
	AutoCreatePortsFromConfig *bool                  `json:"autoCreatePortsFromConfig"`
	Metadata                  *json.RawMessage       `json:"metadata"`
}

// PipelinesParams request payload for querying pipelines.
type PipelinesParams struct {
	Last         *uint
	Before       *string
	Name         *string
	Tags         *string
	ConfigFormat *ConfigFormat
}

// PipelineParams request payload for querying a single pipeline.
type PipelineParams struct {
	ConfigFormat *ConfigFormat
}

// UpdatedPipeline response payload after updating a pipeline successfully.
type UpdatedPipeline struct {
	AddedPorts   []PipelinePort `json:"addedPorts"`
	RemovedPorts []PipelinePort `json:"removedPorts"`
}
