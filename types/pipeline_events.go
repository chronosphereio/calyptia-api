package types

import "time"

type PipelineEventSystem string

const (
	PipelineEventSystemDeployment PipelineEventSystem = "k8s:deployment"
	PipelineEventSystemPod        PipelineEventSystem = "k8s:pod"
)

// PipelineEvent model.
type PipelineEvent struct {
	PipelineID string              `json:"pipeline_id" yaml:"pipeline_id"`
	System     PipelineEventSystem `json:"system" yaml:"system"`
	Status     PipelineStatusKind  `json:"kind" yaml:"kind"`
	Reason     string              `json:"reason" yaml:"reason"`
	Message    string              `json:"message" yaml:"message"`
	CreatedAt  time.Time           `json:"created_at" yaml:"created_at"`
}

// CreatePipelineEvents is the format we submit the events in.
type CreatePipelineEvents []PipelineEvent

// CreatePipelineEventsResponse is the response we get when creating new pipeline
// events.
type CreatePipelineEventsResponse struct {
	Status  string `json:"status" yaml:"status"`
	Message string `json:"message" yaml:"message"`
}
