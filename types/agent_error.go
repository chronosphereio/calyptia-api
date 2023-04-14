package types

import "time"

// AgentError represent an error that occurred on an agent.
// It can be an invalid configuration, or a malfunctioning plugin.
// These come directly from agents.
// Any agent that has errors since the last time it was updated
// will be marked as "errored".
// Errors can be dismissed by the user.
type AgentError struct {
	ID            string     `json:"id" yaml:"id" db:"id"`
	AgentID       string     `json:"agentID" yaml:"agentID" db:"agent_id"`
	Error         string     `json:"error" yaml:"error" db:"error"`
	DismissedAt   *time.Time `json:"dismissedAt" yaml:"dismissedAt" db:"dismissed_at"`
	DismissReason *string    `json:"dismissReason" yaml:"dismissReason" db:"dismiss_reason"`
	CreatedAt     time.Time  `json:"createdAt" yaml:"createdAt" db:"created_at"`
}

type AgentErrors struct {
	Items     []AgentError `json:"items" yaml:"items"`
	EndCursor *string      `json:"endCursor" yaml:"endCursor"`
}

type ListAgentErrors struct {
	AgentID   *string
	FleetID   *string
	Dismissed *bool
	Last      *uint
	Before    *string
}

type CreateAgentError struct {
	AgentID string `json:"-"`
	Error   string `json:"error"`
}

type CreatedAgentError struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

type DismissAgentError struct {
	AgentErrorID string `json:"-"`
	Reason       string `json:"reason"`
}

type DismissedAgentError struct {
	DismissedAt time.Time `json:"dismissedAt" yaml:"dismissedAt"`
}
