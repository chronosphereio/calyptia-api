// Package types contains the required base types used for both client and server side for Calyptia Cloud API.
package types

import (
	"encoding/json"
	"time"
)

const (
	ErrInvalidAgentToken             = UnauthenticatedError("invalid agent token")
	ErrInvalidAgentID                = InvalidArgumentError("invalid agent ID")
	ErrInvalidAgentName              = InvalidArgumentError("invalid agent name")
	ErrInvalidMachineID              = InvalidArgumentError("invalid machine ID")
	ErrInvalidAgentType              = InvalidArgumentError("invalid agent type")
	ErrInvalidAgentVersion           = InvalidArgumentError("invalid agent version")
	ErrUnsupportedAgentVersion       = InvalidArgumentError("unsupported agent version")
	ErrInvalidAgentEdition           = InvalidArgumentError("invalid agent edition")
	ErrInvalidAgentFlags             = InvalidArgumentError("invalid agent flags")
	ErrInvalidAgentFlag              = InvalidArgumentError("invalid agent flag")
	ErrInvalidAgentRawConfig         = InvalidArgumentError("invalid agent raw config")
	ErrInvalidAgentMetadata          = InvalidArgumentError("invalid agent metadata")
	ErrUpdateAgentVersionForbidden   = PermissionDeniedError("update agent version forbidden")
	ErrUpdateAgentEditionForbidden   = PermissionDeniedError("update agent edition forbidden")
	ErrUpdateAgentFlagsForbidden     = PermissionDeniedError("update agent flags forbidden")
	ErrUpdateAgentRawConfigForbidden = PermissionDeniedError("update agent raw config forbidden")
	ErrUpdateAgentMetadataForbidden  = PermissionDeniedError("update agent metadata forbidden")
	ErrAgentGone                     = GoneError("agent gone")
	ErrMachineIDConflict             = ConflictError("machine ID conflict")
	ErrAgentNotFound                 = NotFoundError("agent not found")
)

// Agent model.
type Agent struct {
	ID                  string           `json:"id" yaml:"id"`
	Token               string           `json:"token" yaml:"token"`
	Name                string           `json:"name" yaml:"name"`
	MachineID           string           `json:"machineID" yaml:"machineID"`
	Type                AgentType        `json:"type" yaml:"type"`
	Version             string           `json:"version" yaml:"version"`
	Edition             AgentEdition     `json:"edition" yaml:"edition"`
	Flags               []string         `json:"flags" yaml:"flags"`
	RawConfig           string           `json:"rawConfig" yaml:"rawConfig"`
	Metadata            *json.RawMessage `json:"metadata" yaml:"metadata"`
	FirstMetricsAddedAt time.Time        `json:"firstMetricsAddedAt" yaml:"firstMetricsAddedAt"`
	LastMetricsAddedAt  time.Time        `json:"lastMetricsAddedAt" yaml:"lastMetricsAddedAt"`
	MetricsCount        uint64           `json:"metricsCount" yaml:"metricsCount"`
	CreatedAt           time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

// AgentType declares the fluent agent type (fluentbit/fluentd).
type AgentType string

const (
	// AgentTypeFluentBit fluentbit agent type.
	AgentTypeFluentBit AgentType = "fluentbit"
	// AgentTypeFluentd fluentd agent type.
	AgentTypeFluentd AgentType = "fluentd"
)

// AgentEdition declares the fluent agent edition (community/enterprise).
type AgentEdition string

const (
	// AgentEditionCommunity fluent community edition.
	AgentEditionCommunity AgentEdition = "community"
	// AgentEditionEnterprise fluent enterprise edition.
	AgentEditionEnterprise AgentEdition = "enterprise"
)

// RegisterAgent request payload for registering a new agent.
type RegisterAgent struct {
	Name      string           `json:"name"`
	MachineID string           `json:"machineID"`
	Type      AgentType        `json:"type"`
	Version   string           `json:"version"`
	Edition   AgentEdition     `json:"edition"`
	Flags     []string         `json:"flags"`
	RawConfig string           `json:"rawConfig"`
	Metadata  *json.RawMessage `json:"metadata"`
}

// RegisteredAgent response payload after registering an agent successfully.
type RegisteredAgent struct {
	ID        string    `json:"id"`
	Token     string    `json:"token"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// AgentsParams request payload for querying agents.
type AgentsParams struct {
	Last *uint64
	Name *string
}

// UpdateAgent request payload for updating an agent.
type UpdateAgent struct {
	Name      *string          `json:"name"`
	Version   *string          `json:"version"`
	Edition   *AgentEdition    `json:"edition"`
	Flags     *[]string        `json:"flags"`
	RawConfig *string          `json:"rawConfig"`
	Metadata  *json.RawMessage `json:"metadata"`
}
