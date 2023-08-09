package types

import "time"

type Created struct {
	ID        string    `json:"id" yaml:"id" db:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
}

type Updated struct {
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type UpdatedWithID struct {
	ID        string    `json:"id" yaml:"id" db:"id"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type Deleted struct {
	Deleted   bool       `json:"deleted" yaml:"deleted" db:"deleted"`
	DeletedAt *time.Time `json:"deletedAt" yaml:"deletedAt" db:"deleted_at"`
}
