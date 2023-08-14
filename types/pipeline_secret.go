package types

import "time"

// PipelineSecret model.
type PipelineSecret struct {
	ID        string    `json:"id" yaml:"id"`
	Key       string    `json:"key" yaml:"key"`
	Value     []byte    `json:"value" yaml:"value"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// PipelineSecrets paginated list.
type PipelineSecrets struct {
	Items     []PipelineSecret
	EndCursor *string
}

// CreatePipelineSecret request payload for creating a new pipeline secret.
type CreatePipelineSecret struct {
	Key   string `json:"key"`
	Value []byte `json:"value"`
}

// PipelineSecretsParams request payload for querying the pipeline secrets.
type PipelineSecretsParams struct {
	Last   *uint
	Before *string
}

// UpdatePipelineSecret request payload for updating a pipeline secret.
type UpdatePipelineSecret struct {
	Key   *string `json:"name"`
	Value *[]byte `json:"value"`
}
