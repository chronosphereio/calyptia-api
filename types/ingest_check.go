package types

import (
	"time"
)

const (
	// IngestCheckDefaultTimeout default timeout for an ingest check.
	IngestCheckDefaultTimeout = CoreInstanceNextPingTimeout
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

// CreateIngestCheck request payload for creating a core_instance ingestion check.
type CreateIngestCheck struct {
	Status          CheckStatus    `json:"status"`
	Retries         uint           `json:"retries"`
	ConfigSectionID string         `json:"configSectionID"`
	Logs            []FluentBitLog `json:"logs"`
}

// IngestChecksParams request payload for querying the core_instance ingestion checks.
type IngestChecksParams struct {
	Last   *uint
	Before *string
}

// IngestChecks paginated list.
type IngestChecks struct {
	Items     []IngestCheck `json:"items" yaml:"items"`
	EndCursor *string       `json:"endCursor" yaml:"endCursor"`
}

// UpdateIngestCheck request payload for updating a core_instance ingestion check.
type UpdateIngestCheck struct {
	Status *CheckStatus `json:"status"`

	retries         *uint
	configSectionID *string
	logs            *[]FluentBitLog
}

func (in *UpdateIngestCheck) SetRetries(retries uint) {
	in.retries = &retries
}

func (in *UpdateIngestCheck) SetConfigSectionID(configSectionID string) {
	in.configSectionID = &configSectionID
}

func (in *UpdateIngestCheck) SetLogs(logs []FluentBitLog) {
	in.logs = &logs
}

func (in UpdateIngestCheck) Retries() *uint {
	return in.retries
}

func (in UpdateIngestCheck) ConfigSectionID() *string {
	return in.configSectionID
}

func (in UpdateIngestCheck) Logs() *[]FluentBitLog {
	return in.logs
}
