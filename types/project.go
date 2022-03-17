package types

import "time"

// Project model.
type Project struct {
	ID               string    `json:"id" yaml:"id"`
	Name             string    `json:"name" yaml:"name"`
	MembersCount     uint64    `json:"membersCount" yaml:"membersCount"`
	AgentsCount      uint64    `json:"agentsCount" yaml:"agentsCount"`
	AggregatorsCount uint64    `json:"aggregatorsCount" yaml:"aggregatorsCount"`
	CreatedAt        time.Time `json:"createdAt" yaml:"createdAt"`

	Quotas ProjectQuotas `json:"quotas" yaml:"quotas"`

	Membership *Membership `json:"membership" yaml:"membership"`
}

type ProjectQuotas struct {
	MaxAgents           uint64 `json:"maxAgents" yaml:"maxAgents"`
	MaxAggregators      uint64 `json:"maxAggregators" yaml:"maxAggregators"`
	MaxPipelines        uint64 `json:"maxPipelines" yaml:"maxPipelines"`
	MaxPipelineFiles    uint64 `json:"maxPipelineFiles" yaml:"maxPipelineFiles"`
	MaxPipelineSecrets  uint64 `json:"maxPipelineSecrets" yaml:"maxPipelineSecrets"`
	MaxPipelineFileSize uint64 `json:"maxPipelineFileSize" yaml:"maxPipelineFileSize"`
}

// Projects paginated list.
type Projects struct {
	Items     []Project
	EndCursor *string
}

// CreateProject request payload for creating a project.
type CreateProject struct {
	Name string `json:"name"`
}

// CreatedProject response payload after creating a project successfully.
type CreatedProject struct {
	ID         string     `json:"id"`
	Token      string     `json:"token"`
	CreatedAt  time.Time  `json:"createdAt"`
	Membership Membership `json:"membership"`
}

// ProjectsParams request payload for querying projects.
type ProjectsParams struct {
	Last   *uint64
	Before *string
	Name   *string
}

// UpdateProject request payload for updating a project.
type UpdateProject struct {
	Name *string `json:"name"`
}
