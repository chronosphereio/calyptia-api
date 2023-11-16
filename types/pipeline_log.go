package types

import "time"

type PipelineLogStatus string

const (
	PipelineLogStatusNew  PipelineLogStatus = "NEW"
	PipelineLogStatusDone PipelineLogStatus = "DONE"
)

var AllValidPipelineLogStatuses = [...]PipelineLogStatus{
	PipelineLogStatusNew,
	PipelineLogStatusDone,
}

type PipelineLog struct {
	ID         string            `json:"id" yaml:"id" db:"id"`
	PipelineID string            `json:"pipelineID" yaml:"pipelineID" db:"pipeline_id"`
	Status     PipelineLogStatus `json:"status" yaml:"status" db:"status"`
	Lines      int               `json:"lines" yaml:"lines" db:"lines"`
	Logs       string            `json:"logs" yaml:"logs" db:"logs"`
	CreatedAt  time.Time         `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt  time.Time         `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type CreatePipelineLog struct {
	PipelineID string `json:"pipelineID"`
	Logs       string `json:"logs"`
	Lines      int    `json:"lines"`
}

type UpdatePipelineLog struct {
	ID     string             `json:"id"`
	Logs   *string            `json:"logs"`
	Lines  *int               `json:"lines"`
	Status *PipelineLogStatus `json:"status"`
}

type ListPipelineLogs struct {
	PipelineID string
	Status     *PipelineLogStatus
	Last       *uint
	Before     *string
}

type PipelineLogs struct {
	Items     []PipelineLog `json:"items" yaml:"items"`
	EndCursor *string       `json:"endCursor" yaml:"endCursor"`
	Count     uint          `json:"count" yaml:"count"`
}
