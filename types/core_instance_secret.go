package types

import "time"

type CoreInstanceSecret struct {
	ID             string    `json:"id" yaml:"id" db:"id"`
	CoreInstanceID string    `json:"coreInstanceID" yaml:"coreInstanceID" db:"core_instance_id"`
	Key            string    `json:"key" yaml:"key" db:"key"`
	Value          []byte    `json:"value" yaml:"value" db:"value"`
	CreatedAt      time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreateCoreInstanceSecret struct {
	CoreInstanceID string `json:"-"`
	Key            string `json:"key"`
	Value          []byte `json:"value"`
}

type ListCoreInstanceSecrets struct {
	CoreInstanceID string
	Last           *uint
	Before         *string
}

type CoreInstanceSecrets struct {
	Items     []CoreInstanceSecret `json:"items" yaml:"items"`
	Count     int                  `json:"count" yaml:"count"`
	EndCursor *string              `json:"endCursor" yaml:"endCursor"`
}

type UpdateCoreInstanceSecret struct {
	ID    string  `json:"-"`
	Key   *string `json:"key"`
	Value *[]byte `json:"value"`
}
