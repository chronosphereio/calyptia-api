package types

import "time"

type CoreInstanceFile struct {
	ID             string    `json:"id" yaml:"id" db:"id"`
	CoreInstanceID string    `json:"coreInstanceID" yaml:"coreInstanceID" db:"core_instance_id"`
	Name           string    `json:"name" yaml:"name" db:"name"`
	Contents       []byte    `json:"contents" yaml:"contents" db:"contents"`
	Encrypted      bool      `json:"encrypted" yaml:"encrypted" db:"encrypted"`
	CreatedAt      time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreateCoreInstanceFile struct {
	CoreInstanceID string `json:"-"`
	Name           string `json:"name"`
	Contents       []byte `json:"contents"`
	Encrypted      bool   `json:"encrypted"`
}

type ListCoreInstanceFiles struct {
	CoreInstanceID string
	Last           *uint
	Before         *string
}

type CoreInstanceFiles struct {
	Items     []CoreInstanceFile `json:"items" yaml:"items"`
	Count     int                `json:"count" yaml:"count"`
	EndCursor *string            `json:"endCursor" yaml:"endCursor"`
}

type UpdateCoreInstanceFile struct {
	ID        string  `json:"-"`
	Name      *string `json:"name"`
	Contents  *[]byte `json:"contents"`
	Encrypted *bool   `json:"encrypted"`
}
