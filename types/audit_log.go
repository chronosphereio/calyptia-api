package types

import "time"

type CreateAuditLog struct {
	Timestamp    time.Time         `json:"timestamp" yaml:"timestamp" db:"timestamp"`
	ProjectID    string            `json:"projectID" yaml:"projectID" db:"projectID"`
	URL          string            `json:"url" yaml:"url" db:"url"`
	Components   map[string]string `json:"components" yaml:"components" db:"components"`
	Action       string            `json:"action" yaml:"action" db:"action"`
	Identity     any               `json:"identity" yaml:"identity" db:"identity"`
	ResponseCode int               `json:"responseCode" yaml:"responseCode" db:"response_code"`
}

type CreatedAuditLog struct {
	ID string `json:"id"`
}

type AuditLog struct {
	ID           string            `json:"id" yaml:"id" db:"id"`
	Timestamp    time.Time         `json:"timestamp" yaml:"timestamp" db:"created_at"`
	URL          string            `json:"url" yaml:"url" db:"url"`
	Components   map[string]string `json:"components" yaml:"components" db:"components"`
	Action       string            `json:"action" yaml:"action" db:"action"`
	Identity     any               `json:"identity" yaml:"identity" db:"identity"`
	ResponseCode int               `json:"responseCode" yaml:"responseCode" db:"response_code"`
}

type AuditLogs struct {
	Items     []AuditLog `json:"items" yaml:"items"`
	EndCursor *string    `json:"endCursor" yaml:"endCursor"`
}

type ListAuditLogs struct {
	// project scope
	ProjectID string

	URL string

	// audit log specific queries
	Action     string
	Components map[string]string
	Identity   map[string]string

	// pagination
	First *uint
	After *string
}
