package types

import (
	"time"
)

// PipelineCheck type for pipeline level checks.
type PipelineCheck struct {
	Check
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// CreatePipelineCheck request payload for creating a core_instance check.
type CreatePipelineCheck Check

// PipelineChecksParams request payload for querying the core_instance checks.
type PipelineChecksParams struct {
	Last   *uint
	Before *string
}

// PipelineChecks paginated list.
type PipelineChecks struct {
	Items     []PipelineCheck `json:"items" yaml:"items"`
	EndCursor *string         `json:"endCursor" yaml:"endCursor"`
}

// UpdatePipelineCheck request payload for updating a core_instance check.
type UpdatePipelineCheck struct {
	Protocol *PipelinePortProtocol `json:"protocol"`
	Status   *CheckStatus          `json:"status"`
	Retries  *uint                 `json:"retries"`
	Host     *string               `json:"host"`
	Port     *uint                 `json:"port"`
}
