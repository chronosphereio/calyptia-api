package types

import "time"

const (
	ErrInvalidPipelinePortID                = InvalidArgumentError("invalid pipeline port ID")
	ErrInvalidPipelinePortProtocol          = InvalidArgumentError("invalid pipeline port protocol")
	ErrInvalidPipelinePortNumber            = InvalidArgumentError("invalid pipeline port number")
	ErrInvalidPipelinePortEndpoint          = InvalidArgumentError("invalid pipeline port endpoint")
	ErrPipelinePortFrontendAlreadyAllocated = ConflictError("pipeline port frontend already allocated")
	ErrUpdatePipelinePortEndpointForbidden  = PermissionDeniedError("update pipeline port endpoint forbidden")
	ErrPipelinePortNotFound                 = NotFoundError("pipeline port not found")
)

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

// CreatedPipelinePort response payload after creating a pipeline port successfully.
type CreatedPipelinePort struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

// PipelinePortsParams request payload for querying the pipeline ports.
type PipelinePortsParams struct {
	Last *uint64
}

// UpdatePipelinePort request payload for updating a pipeline port.
type UpdatePipelinePort struct {
	Protocol     *string `json:"protocol"`
	FrontendPort *uint   `json:"frontendPort"`
	BackendPort  *uint   `json:"backendPort"`
	Endpoint     *string `json:"endpoint"`
}
