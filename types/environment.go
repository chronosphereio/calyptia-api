package types

import "time"

type Environment struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	ProjectID string    `json:"project_id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

type CreateEnvironment struct {
	Name      string `json:"name"`
	ProjectID string `json:"project_id"`
}

type CreatedEnvironment struct {
	ID         string     `json:"id"`
	CreatedAt  time.Time  `json:"createdAt"`
	Membership Membership `json:"membership"`
}

type Environments struct {
	Items     []Environment
	EndCursor *string
}

type EnvironmentsParams struct {
	Last   *uint64
	Before *string
}

type UpdateEnvironment struct {
	Name string `json:"name"`
}
