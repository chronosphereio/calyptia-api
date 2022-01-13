package types

import "time"

// AgentConfig model.
type AgentConfig struct {
	ID        string    `json:"id" yaml:"id"`
	RawConfig string    `json:"rawConfig" yaml:"rawConfig"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// AgentConfigHistoryParams request payload for querying the agent config history.
type AgentConfigHistoryParams struct {
	Last *uint64
}
