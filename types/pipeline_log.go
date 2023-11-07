package types

import "time"

type PipelineLog struct {
	ID         string    `json:"id" yaml:"id" db:"id"`
	PipelineID string    `json:"pipelineID" yaml:"pipelineID" db:"pipeline_id"`
	Logs       string    `json:"logs" yaml:"logs" db:"logs"`
	CreatedAt  time.Time `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreatePipelineLog struct {
	PipelineID string `json:"pipelineID"`
	Logs       string `json:"logs"`
}

type UpdatePipelineLog struct {
	ID   string  `json:"id"`
	Logs *string `json:"logs"`
}

type ListPipelineLogs struct {
	PipelineID string
	Last       *uint
	Before     *string
}

type PipelineLogs struct {
	Items     []PipelineLog `json:"items" yaml:"items"`
	EndCursor *string       `json:"endCursor" yaml:"endCursor"`
	Count     uint          `json:"count" yaml:"count"`
}
