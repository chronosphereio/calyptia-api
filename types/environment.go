package types

import "time"

type Environment struct {
	ID        string    `json:"id" yaml:"id" `
	Name      string    `json:"name" yaml:"name" `
	ProjectID string    `json:"-" yaml:"-"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

type CreateEnvironment struct {
	Name     string `json:"name"`
	internal bool
}

func (in *CreateEnvironment) SetInternal(internal bool) {
	in.internal = internal
}

func (in CreateEnvironment) Internal() bool {
	return in.internal
}

type CreatedEnvironment struct {
	ID         string     `json:"id" yaml:"id"`
	CreatedAt  time.Time  `json:"createdAt" yaml:"createdAt"`
	Membership Membership `json:"membership" yaml:"membership"`
}

type Environments struct {
	Items     []Environment
	EndCursor *string
}

type EnvironmentsParams struct {
	Last   *uint
	Before *string
	Name   *string
}

type UpdateEnvironment struct {
	Name *string `json:"name"`
}
