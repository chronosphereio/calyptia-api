package types

import "time"

// PipelinePort model.
type PipelinePort struct {
	ID           string    `json:"id" yaml:"id"`
	Protocol     string    `json:"protocol" yaml:"protocol"`
	FrontendPort uint      `json:"frontendPort" yaml:"frontendPort"`
	BackendPort  uint      `json:"backendPort" yaml:"backendPort"`
	Endpoint     string    `json:"endpoint" yaml:"endpoint"`
	CreatedAt    time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// CreatePipelinePort request payload for creating a pipeline port.
type CreatePipelinePort struct {
	Protocol     string `json:"protocol"`
	FrontendPort uint   `json:"frontendPort"`
	BackendPort  uint   `json:"backendPort"`
	Endpoint     string `json:"endpoint"`
}

// PipelinePortsParams request payload for querying the pipeline ports.
type PipelinePortsParams struct {
	Last *uint64
}

// UpdatePipelinePortOpts request payload for updating a pipeline port.
type UpdatePipelinePortOpts struct {
	Protocol     *string `json:"protocol"`
	FrontendPort *uint   `json:"frontendPort"`
	BackendPort  *uint   `json:"backendPort"`
	Endpoint     *string `json:"endpoint"`
}
