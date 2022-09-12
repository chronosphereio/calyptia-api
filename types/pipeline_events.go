package types

import "time"

type PipelineEventSystem string

const (
	PipelineEventSystemDeployment PipelineEventSystem = "k8s:deployment"
	PipelineEventSystemPod        PipelineEventSystem = "k8s:pod"
)

// PipelineEvent model.
type PipelineEvent struct {
	PipelineID string              `json:"pipelineID" yaml:"pipelineID"`
	System     PipelineEventSystem `json:"system" yaml:"system"`
	Status     PipelineStatusKind  `json:"kind" yaml:"kind"`
	Reason     string              `json:"reason" yaml:"reason"`
	Message    string              `json:"message" yaml:"message"`
	LoggedAt   time.Time           `json:"loggedAt" yaml:"loggedAt"`
}

// CreatePipelineEvents is the format we submit the events in.
type CreatePipelineEvent struct {
	PipelineEvent
}

// CreatedPipelineEvent is the response we get when creating new pipeline
// events.
type CreatedPipelineEvent struct {
	CreatePipelineEvent
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}
