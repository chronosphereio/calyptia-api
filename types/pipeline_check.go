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

// PipelineChecks paginated list.
type PipelineChecks struct {
	Items     []PipelineCheck `json:"items" yaml:"items"`
	EndCursor *string         `json:"endCursor" yaml:"endCursor"`
}

// PipelineChecksParams request payload for querying the core_instance checks.
type PipelineChecksParams struct {
	Last   *uint
	Before *string
}

// UpdatePipelineCheck request payload for updating a core_instance check.
type UpdatePipelineCheck struct {
	Protocol *PipelinePortProtocol `json:"protocol"`
	Status   *CheckStatus          `json:"status"`
	Retries  *uint                 `json:"retries"`
	Host     *string               `json:"host"`
	Port     *uint                 `json:"port"`
}

// CreatePipelineCheck request payload for creating a core_instance check.
type CreatePipelineCheck Check

// CreatedPipelineCheck response payload after creating an aggregator successfully.
type CreatedPipelineCheck struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
