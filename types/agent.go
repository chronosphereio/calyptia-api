// Package types contains the required base types used for both client and server side for Calyptia Cloud API.
package types

import (
	"encoding/json"
	"time"
)

// Agent model.
type Agent struct {
	ID                  string           `json:"id" yaml:"id"`
	FleetID             *string          `json:"fleetID" yaml:"fleetID"`
	Token               string           `json:"token" yaml:"token"`
	Name                string           `json:"name" yaml:"name"`
	MachineID           string           `json:"machineID" yaml:"machineID"`
	EnvironmentName     string           `json:"environmentName" yaml:"environmentName"`
	Type                AgentType        `json:"type" yaml:"type"`
	Version             string           `json:"version" yaml:"version"`
	Edition             AgentEdition     `json:"edition" yaml:"edition"`
	Flags               []string         `json:"flags" yaml:"flags"`
	RawConfig           string           `json:"rawConfig" yaml:"rawConfig"`
	Metadata            *json.RawMessage `json:"metadata" yaml:"metadata"`
	Tags                []string         `json:"tags" yaml:"tags"`
	FirstMetricsAddedAt *time.Time       `json:"firstMetricsAddedAt" yaml:"firstMetricsAddedAt"`
	LastMetricsAddedAt  *time.Time       `json:"lastMetricsAddedAt" yaml:"lastMetricsAddedAt"`
	MetricsCount        uint             `json:"metricsCount" yaml:"metricsCount"`
	CreatedAt           time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

// Agents paginated list.
type Agents struct {
	Items     []Agent
	EndCursor *string
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
	EnvironmentID string           `json:"environmentID"`
	FleetID       *string          `json:"fleetID"`
	Name          string           `json:"name"`
	MachineID     string           `json:"machineID"`
	Type          AgentType        `json:"type"`
	Version       string           `json:"version"`
	Edition       AgentEdition     `json:"edition"`
	Flags         []string         `json:"flags"`
	RawConfig     string           `json:"rawConfig"`
	Metadata      *json.RawMessage `json:"metadata"`
	Tags          []string         `json:"tags"`

	id         string
	signingKey []byte
	token      string
}

func (in *RegisterAgent) SetID(id string) {
	in.id = id
}

func (in *RegisterAgent) SetSigningKey(signingKey []byte) {
	in.signingKey = signingKey
}

func (in *RegisterAgent) SetToken(token string) {
	in.token = token
}

func (in RegisterAgent) ID() string {
	return in.id
}

func (in RegisterAgent) SigningKey() []byte {
	return in.signingKey
}

func (in RegisterAgent) Token() string {
	return in.token
}

// RegisteredAgent response payload after registering an agent successfully.
type RegisteredAgent struct {
	ID              string    `json:"id"`
	Token           string    `json:"token"`
	Name            string    `json:"name"`
	CreatedAt       time.Time `json:"createdAt"`
	EnvironmentName string    `json:"environmentName"`
}

// AgentsParams request payload for querying agents.
type AgentsParams struct {
	Last          *uint
	Before        *string
	Name          *string
	FleetID       *string
	EnvironmentID *string
	TagsQuery     *string
	tags          []string
}

func (in *AgentsParams) SetTags(tags []string) {
	in.tags = tags
}

func (in AgentsParams) Tags() []string {
	return in.tags
}

// UpdateAgent request payload for updating an agent.
type UpdateAgent struct {
	FleetID       *string          `json:"fleetID"`
	EnvironmentID *string          `json:"environmentID"`
	Name          *string          `json:"name"`
	Version       *string          `json:"version"`
	Edition       *AgentEdition    `json:"edition"`
	Flags         *[]string        `json:"flags"`
	RawConfig     *string          `json:"rawConfig"`
	Metadata      *json.RawMessage `json:"metadata"`

	firstMetricsAddedAt *time.Time
	lastMetricsAddedAt  *time.Time
	newMetricsCount     *uint
}

func (in *UpdateAgent) SetFirstMetricsAddedAt(firstMetricsAddedAt *time.Time) {
	in.firstMetricsAddedAt = firstMetricsAddedAt
}

func (in *UpdateAgent) SetLastMetricsAddedAt(lastMetricsAddedAt *time.Time) {
	in.lastMetricsAddedAt = lastMetricsAddedAt
}

func (in *UpdateAgent) SetNewMetricsCount(newMetricsCount *uint) {
	in.newMetricsCount = newMetricsCount
}

func (in UpdateAgent) FirstMetricsAddedAt() *time.Time {
	return in.firstMetricsAddedAt
}

func (in UpdateAgent) LastMetricsAddedAt() *time.Time {
	return in.lastMetricsAddedAt
}

func (in UpdateAgent) NewMetricsCount() *uint {
	return in.newMetricsCount
}
