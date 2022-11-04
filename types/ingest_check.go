package types

import (
	"time"
)

const (
	// IngestCheckDefaultTimeout default timeout for an ingest check.
	IngestCheckDefaultTimeout = AggregatorNextPingTimeout
)

// IngestCheck type for core_instance ingestion checks.
type IngestCheck struct {
	ID              string      `json:"id" yaml:"id"`
	ConfigSectionID string      `json:"-" yaml:"-"`
	Config          string      `json:"config" yaml:"config"`
	Status          CheckStatus `json:"status" yaml:"status"`
	Retries         uint        `json:"retries" yaml:"retries"`
	CreatedAt       time.Time   `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time   `json:"updatedAt" yaml:"updatedAt"`
}

// IngestChecks paginated list.
type IngestChecks struct {
	Items     []IngestCheck `json:"items" yaml:"items"`
	EndCursor *string       `json:"endCursor" yaml:"endCursor"`
}

// IngestChecksParams request payload for querying the core_instance ingestion checks.
type IngestChecksParams struct {
	Last   *uint
	Before *string
}

// UpdateIngestCheck request payload for updating a core_instance ingestion check.
type UpdateIngestCheck struct {
	Status *CheckStatus `json:"status"`
}

// CreateIngestCheck request payload for creating a core_instance ingestion check.
type CreateIngestCheck struct {
	Status          CheckStatus    `json:"status"`
	Retries         uint           `json:"retries"`
	ConfigSectionID string         `json:"configSectionID"`
	Logs            []FluentBitLog `json:"logs"`
}

// CreatedIngestCheck response payload after creating a core ingestion check successfully.
type CreatedIngestCheck struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
