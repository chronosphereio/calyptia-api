package types

import "time"

// PipelineFile model.
type PipelineFile struct {
	ID              string    `json:"id" yaml:"id"`
	Name            string    `json:"name" yaml:"name"`
	Contents        []byte    `json:"contents" yaml:"contents"`
	Encrypted       bool      `json:"encrypted" yaml:"encrypted"`
	ProcesingRuleID *string   `json:"procesingRuleID" yaml:"procesingRuleID"`
	CreatedAt       time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// PipelineFiles paginated list.
type PipelineFiles struct {
	Items     []PipelineFile
	EndCursor *string
}

// CreatePipelineFile request payload for creating a new pipeline file.
type CreatePipelineFile struct {
	Name      string `json:"name"`
	Contents  []byte `json:"contents"`
	Encrypted bool   `json:"encrypted"`
}

// CreatedPipelineFile response payload after creating a pipeline file successfully.
type CreatedPipelineFile struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

// PipelineFilesParams request payload for querying the pipeline files.
type PipelineFilesParams struct {
	Last   *uint
	Before *string
}

// UpdatePipelineFile request payload for updating a pipeline file.
type UpdatePipelineFile struct {
	Name      *string `json:"name"`
	Contents  *[]byte `json:"contents"`
	Encrypted *bool   `json:"encrypted"`
}
