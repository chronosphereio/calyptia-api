package types

import "time"

type SAMLMapping struct {
	ID         string    `json:"id" yaml:"id" db:"id"`
	ProjectID  string    `json:"projectID" yaml:"projectID" db:"project_id"`
	ClaimKey   string    `json:"claimKey" yaml:"claimKey" db:"claim_key"`
	ClaimValue string    `json:"claimValue" yaml:"claimValue" db:"claim_value"`
	Role       Role      `json:"role" yaml:"role" db:"role"`
	CreatedAt  time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreateSAMLMapping struct {
	ProjectID  string `json:"projectID"`
	ClaimKey   string `json:"claimKey"`
	ClaimValue string `json:"claimValue"`
	Role       Role   `json:"role"`
}

type ListSAMLMappings struct {
	ProjectID string
	Last      *uint
	Before    *string
}

type SAMLMappings struct {
	Items     []SAMLMapping `json:"items" yaml:"items"`
	EndCursor *string       `json:"endCursor" yaml:"endCursor"`
	Count     int           `json:"count" yaml:"count"`
}

type UpdateSAMLMapping struct {
	ID         string  `json:"-"`
	ClaimKey   *string `json:"claimKey"`
	ClaimValue *string `json:"claimValue"`
	Role       *Role   `json:"role"`
}
