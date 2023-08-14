package types

import (
	"time"
)

// CoreInstanceCheck type for pipeline level checks.
type CoreInstanceCheck struct {
	Check
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// CoreInstanceChecks paginated list.
type CoreInstanceChecks struct {
	Items     []CoreInstanceCheck `json:"items" yaml:"items"`
	EndCursor *string             `json:"endCursor" yaml:"endCursor"`
}

// CoreInstanceChecksParams request payload for querying the core_instance checks.
type CoreInstanceChecksParams struct {
	Last   *uint
	Before *string
}

// UpdateCoreInstanceCheck request payload for updating a core_instance check.
type UpdateCoreInstanceCheck struct {
	Protocol *PipelinePortProtocol `json:"protocol"`
	Status   *CheckStatus          `json:"status"`
	Host     *string               `json:"host"`
	Retries  *uint                 `json:"retries"`
	Port     *uint                 `json:"port"`
}

// CreateCoreInstanceCheck request payload for creating a core_instance check.
type CreateCoreInstanceCheck = Check
