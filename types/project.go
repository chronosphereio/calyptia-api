package types

import "time"

// Project model.
type Project struct {
	ID                 string    `json:"id" yaml:"id"`
	Name               string    `json:"name" yaml:"name"`
	MembersCount       uint      `json:"membersCount" yaml:"membersCount"`
	AgentsCount        uint      `json:"agentsCount" yaml:"agentsCount"`
	CoreInstancesCount uint      `json:"aggregatorsCount" yaml:"aggregatorsCount"`
	CreatedAt          time.Time `json:"createdAt" yaml:"createdAt"`

	Membership *Membership `json:"membership" yaml:"membership"`
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
	Last   *uint
	Before *string
	Name   *string
}

// UpdateProject request payload for updating a project.
type UpdateProject struct {
	Name *string `json:"name"`
}
