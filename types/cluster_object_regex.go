package types

import "time"

type ClusterObjectRegex struct {
	ID          string    `json:"id" yaml:"id" db:"id"`
	PipelineID  string    `json:"pipelineID" yaml:"pipelineID" db:"pipeline_id"`
	Regex       string    `json:"regex" yaml:"regex" db:"regex"`
	Description string    `json:"description" yaml:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`

	ClusterObjects *[]ClusterObject `json:"clusterObjects,omitempty" yaml:"clusterObjects,omitempty" db:"-"`
}

type CreateClusterObjectRegex struct {
	PipelineID  string `json:"-"`
	Regex       string `json:"regex"`
	Description string `json:"description"`
}

type ListClusterObjectRegexes struct {
	PipelineID string
	Last       *uint
	Before     *string
}

type ClusterObjectRegexes struct {
	Items     []ClusterObjectRegex `json:"items"`
	EndCursor *string              `json:"endCursor"`
}

type UpdateClusterObjectRegex struct {
	RegexID     string  `json:"-"`
	Regex       *string `json:"regex"`
	Description *string `json:"description"`
}
