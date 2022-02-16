package types

import "time"

const (
	ErrInvalidProjectID   = InvalidArgumentError("invalid project ID")
	ErrInvalidProjectName = InvalidArgumentError("invalid project name")
	ErrProjectGone        = GoneError("project gone")
	ErrProjectNotFound    = NotFoundError("project not found")
)

// Project model.
type Project struct {
	ID               string    `json:"id" yaml:"id"`
	Name             string    `json:"name" yaml:"name"`
	MembersCount     uint64    `json:"membersCount" yaml:"membersCount"`
	AgentsCount      uint64    `json:"agentsCount" yaml:"agentsCount"`
	AggregatorsCount uint64    `json:"aggregatorsCount" yaml:"aggregatorsCount"`
	CreatedAt        time.Time `json:"createdAt" yaml:"createdAt"`

	Membership *Membership `json:"membership" yaml:"membership"`
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
	Last *uint64
	Name *string
}

// UpdateProject request payload for updating a project.
type UpdateProject struct {
	Name *string `json:"name"`
}
