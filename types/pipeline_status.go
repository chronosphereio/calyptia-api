package types

import "time"

type PipelineEventSystem string

const (
	PipelineEventSystemDeployment PipelineEventSystem = "k8s:deployment"
	PipelineEventSystemPod        PipelineEventSystem = "k8s:pod"
)

// PipelineEvent model.
type PipelineEvent struct {
	System   PipelineEventSystem `json:"system" yaml:"system"`
	Reason   string              `json:"reason" yaml:"reason"`
	Message  string              `json:"message" yaml:"message"`
	LoggedAt time.Time           `json:"loggedAt" yaml:"loggedAt"`
}

// PipelineStatus model.
type PipelineStatus struct {
	ID        string             `json:"id" yaml:"id"`
	Config    PipelineConfig     `json:"config" yaml:"config"`
	Status    PipelineStatusKind `json:"status" yaml:"status"`
	Events    []PipelineEvent    `json:"events,omitempty" yaml:"events,omitempty"`
	CreatedAt time.Time          `json:"createdAt" yaml:"createdAt"`
}

// PipelineStatusHistory paginated list.
type PipelineStatusHistory struct {
	Items     []PipelineStatus
	EndCursor *string
}

// PipelineStatusKind enum.
type PipelineStatusKind string

const (
	// PipelineStatusNew is the default status right after a pipeline is created.
	PipelineStatusNew PipelineStatusKind = "NEW"
	// PipelineStatusFailed is the status when a pipeline fails.
	PipelineStatusFailed PipelineStatusKind = "FAILED"
	// PipelineStatusStarting is the status of a starting pipeline.
	PipelineStatusStarting PipelineStatusKind = "STARTING"
	// PipelineStatusStarted is the status of a started pipeline.
	PipelineStatusStarted PipelineStatusKind = "STARTED"
	// PipelineStatusScaling is the status of a pipeline while scaling up/down.
	PipelineStatusScaling PipelineStatusKind = "SCALING"
)

// PipelineStatusHistoryParams request payload for querying the pipeline status history.
type PipelineStatusHistoryParams struct {
	Last         *uint
	Before       *string
	Status       *PipelineStatusKind
	ConfigFormat *ConfigFormat
}
