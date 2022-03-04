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

// PipelinePorts paginated list.
type PipelinePorts struct {
	Items     []PipelinePort
	EndCursor *string
}

// CreatePipelinePort request payload for creating a pipeline port.
type CreatePipelinePort struct {
	Protocol     string `json:"protocol"`
	FrontendPort uint   `json:"frontendPort"`
	BackendPort  uint   `json:"backendPort"`
	Endpoint     string `json:"endpoint"`
}

// CreatedPipelinePort response payload after creating a pipeline port successfully.
type CreatedPipelinePort struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

// PipelinePortsParams request payload for querying the pipeline ports.
type PipelinePortsParams struct {
	Last   *uint64
	Before *string
}

// UpdatePipelinePort request payload for updating a pipeline port.
type UpdatePipelinePort struct {
	Protocol     *string `json:"protocol"`
	FrontendPort *uint   `json:"frontendPort"`
	BackendPort  *uint   `json:"backendPort"`
	Endpoint     *string `json:"endpoint"`
}
