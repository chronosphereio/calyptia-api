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
	Status              AgentStatus      `json:"status" yaml:"status"`
	Version             string           `json:"version" yaml:"version"`
	Edition             AgentEdition     `json:"edition" yaml:"edition"`
	OperatingSystem     AgentOS          `json:"os" yaml:"os"`
	Architecture        AgentArch        `json:"arch" yaml:"arch"`
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

type AgentStatus string

const (
	AgentStatusHealthy   AgentStatus = "healthy"
	AgentStatusUnhealthy AgentStatus = "unhealthy"
)

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

// AgentOS is set to the operating system the agent is running on.
type AgentOS string

const (
	// AgentOSUnknown is the default value for the operating system.
	AgentOSUnknown AgentOS = "unknown"
	// AgentOSWindows is for Win32 machines.
	AgentOSWindows AgentOS = "windows"
	// AgentOSMacOS is for macOS machines.
	AgentOSMacOS AgentOS = "macos"
	// AgentOSLinux is for Win32 machines.
	AgentOSLinux AgentOS = "linux"
	// AgentOSFreeBSD is for FreeBSD machines.
	AgentOSFreeBSD AgentOS = "freebsd"
	// AgentOSOpenBSD is for OpenBSD machines.
	AgentOSOpenBSD AgentOS = "openbsd"
	// AgentOSNetBSD is for NetBSD machines.
	AgentOSNetBSD AgentOS = "netbsd"
)

// AgentArch is set to the architecture an agent is running on.
type AgentArch string

const (
	// AgentArchUnknown is the default value for architecture.
	AgentArchUnknown AgentArch = "unknown"
	// AgentArchX86 is for intel i686 machines.
	AgentArchX86 AgentArch = "x86"
	// AgentArchX86_64 is for intel x86_64 machines.
	AgentArchX86_64 AgentArch = "x86_64"
	// AgentArchArm is for arm machines.
	AgentArchArm AgentArch = "arm"
	// AgentArchArm64 is for arm64 machines.
	AgentArchArm64 AgentArch = "arm64"
)

// RegisterAgent request payload for registering a new agent.
type RegisterAgent struct {
	EnvironmentID   string           `json:"environmentID"`
	FleetID         *string          `json:"fleetID"`
	Name            string           `json:"name"`
	MachineID       string           `json:"machineID"`
	Type            AgentType        `json:"type"`
	Version         string           `json:"version"`
	Edition         AgentEdition     `json:"edition"`
	OperatingSystem *AgentOS         `json:"os"`
	Architecture    *AgentArch       `json:"arch"`
	Flags           []string         `json:"flags"`
	RawConfig       string           `json:"rawConfig"`
	Metadata        *json.RawMessage `json:"metadata"`
	Tags            []string         `json:"tags"`

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
	FleetID         *string          `json:"fleetID"`
	EnvironmentID   *string          `json:"environmentID"`
	Name            *string          `json:"name"`
	Version         *string          `json:"version"`
	Edition         *AgentEdition    `json:"edition"`
	OperatingSystem *AgentOS         `json:"os"`
	Architecture    *AgentArch       `json:"arch"`
	Flags           *[]string        `json:"flags"`
	RawConfig       *string          `json:"rawConfig"`
	Metadata        *json.RawMessage `json:"metadata"`

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
