package types

import (
	"encoding/json"
	"time"
)

type PipelineKind string

const (
	PipelineKindDaemonSet  PipelineKind = "daemonSet"
	PipelineKindDeployment PipelineKind = "deployment"
)

// AllPipelineKindTypes all valid pipeline kinds.
var AllPipelineKindTypes = [...]PipelineKind{
	PipelineKindDaemonSet,
	PipelineKindDeployment,
}

// Pipeline model.
type Pipeline struct {
	ID                           string           `json:"id" yaml:"id"`
	Name                         string           `json:"name" yaml:"name"`
	Kind                         PipelineKind     `json:"kind" yaml:"kind"`
	Config                       PipelineConfig   `json:"config" yaml:"config"`
	ConfigSections               []ConfigSection  `json:"configSections" yaml:"configSections"`
	Image                        *string          `json:"image" yaml:"image"`
	Status                       PipelineStatus   `json:"status" yaml:"status"`
	ResourceProfile              ResourceProfile  `json:"resourceProfile" yaml:"resourceProfile"`
	TracingEnabled               bool             `json:"tracingEnabled" yaml:"tracingEnabled"`
	WaitForChecksBeforeDeploying bool             `json:"waitForChecksBeforeDeploying" yaml:"waitForChecksBeforeDeploying"`
	ReplicasCount                uint             `json:"replicasCount" yaml:"replicasCount"`
	Tags                         []string         `json:"tags" yaml:"tags"`
	Metadata                     *json.RawMessage `json:"metadata" yaml:"metadata"`
	ChecksTotal                  uint             `json:"checksTotal" yaml:"checksTotal"`
	ChecksOK                     uint             `json:"checksOK" yaml:"checksOK"`
	ChecksRunning                uint             `json:"checksRunning" yaml:"checksRunning"`
	CreatedAt                    time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt                    time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

// Pipelines paginated list.
type Pipelines struct {
	Items     []Pipeline
	EndCursor *string
}

// CreatePipeline request payload for creating a new pipeline.
type CreatePipeline struct {
	Name                 string                 `json:"name"`
	Kind                 PipelineKind           `json:"kind"`
	ReplicasCount        uint                   `json:"replicasCount"`
	RawConfig            string                 `json:"rawConfig"`
	ConfigFormat         ConfigFormat           `json:"configFormat"`
	Secrets              []CreatePipelineSecret `json:"secrets"`
	Files                []CreatePipelineFile   `json:"files"`
	ResourceProfileName  string                 `json:"resourceProfile"`
	Image                *string                `json:"image"`
	SkipConfigValidation bool                   `json:"skipConfigValidation"`
	Metadata             *json.RawMessage       `json:"metadata"`
	Tags                 []string               `json:"tags"`

	// Deprecated: in favor of NoAutoCreateEndpointsFromConfig
	AutoCreatePortsFromConfig bool `json:"autoCreatePortsFromConfig"`
	// Deprecated: in favor of NoAutoCreateChecksFromConfig
	AutoCreateChecksFromConfig bool `json:"autoCreateChecksFromConfig"`

	// no automatically create endpoints from config
	NoAutoCreateEndpointsFromConfig bool `json:"noAutoCreateEndpointsFromConfig"`
	// no automatically create checks based on the output configuration.
	NoAutoCreateChecksFromConfig bool `json:"noAutoCreateChecksFromConfig"`

	// WaitForChecksBeforeDeploying is a conditional variable that defines behavior on the
	// pipeline deployment
	//
	// If set to true:
	//
	// If all checks associated with the pipeline run successfully, the status of
	// the pipeline will be switched to CHECKS_OK and the deployment will be executed.
	//
	// If any of the checks associated with the pipeline fails, the status of
	// the pipeline will be switched to CHECKS_FAILED and the deployment of the pipeline
	// will be blocked.
	//
	// If set to false (default):
	//
	// If all checks associated with the pipeline run successfully, the status of
	// the pipeline will be switched to CHECKS_OK and the deployment will be executed.
	//
	// If any of the checks associated with the pipeline fails, the status of
	// the pipeline will be switched to CHECKS_FAILED and the deployment of the pipeline
	// will be executed.
	WaitForChecksBeforeDeploying bool `json:"waitForChecksBeforeDeploying"`
}

// CreatedPipeline response payload after creating a pipeline successfully.
type CreatedPipeline struct {
	ID                           string           `json:"id"`
	Name                         string           `json:"name"`
	Kind                         PipelineKind     `json:"kind"`
	Config                       PipelineConfig   `json:"config"`
	Secrets                      []PipelineSecret `json:"secrets"`
	Files                        []PipelineFile   `json:"files"`
	Status                       PipelineStatus   `json:"status"`
	ResourceProfile              ResourceProfile  `json:"resourceProfile"`
	Checks                       []PipelineCheck  `json:"checks"`
	ReplicasCount                uint             `json:"replicasCount"`
	WaitForChecksBeforeDeploying bool             `json:"waitForChecksBeforeDeploying"`
	CreatedAt                    time.Time        `json:"createdAt"`
}

// UpdatePipeline request payload for updating a pipeline.
type UpdatePipeline struct {
	Name            *string             `json:"name"`
	Kind            *PipelineKind       `json:"kind"`
	Status          *PipelineStatusKind `json:"status"`
	ConfigFormat    *ConfigFormat       `json:"configFormat"`
	ReplicasCount   *uint               `json:"replicasCount"`
	RawConfig       *string             `json:"rawConfig"`
	ResourceProfile *string             `json:"resourceProfile"`
	Image           *string             `json:"image"`
	// Deprecated: in favor of NoAutoCreateEndpointsFromConfig
	AutoCreatePortsFromConfig *bool `json:"autoCreatePortsFromConfig"`
	// Deprecated: in favor of NoAutoCreateChecksFromConfig
	AutoCreateChecksFromConfig *bool `json:"autoCreateChecksFromConfig"`

	// no automatically create endpoints from config
	NoAutoCreateEndpointsFromConfig bool `json:"noAutoCreateEndpointsFromConfig"`
	// no automatically create checks based on the output configuration.
	NoAutoCreateChecksFromConfig bool `json:"noAutoCreateChecksFromConfig"`

	// this defines behavior; await for checks to complete before reporting the status back.
	WaitForChecksBeforeDeploying *bool `json:"waitForChecksBeforeDeploying"`

	SkipConfigValidation bool                   `json:"skipConfigValidation"`
	Metadata             *json.RawMessage       `json:"metadata"`
	Secrets              []UpdatePipelineSecret `json:"secrets"`
	Files                []UpdatePipelineFile   `json:"files"`
	Events               []PipelineEvent        `json:"events"`
}

// PipelinesParams request payload for querying pipelines.
type PipelinesParams struct {
	Last                     *uint
	Before                   *string
	Name                     *string
	Tags                     *string
	ConfigFormat             *ConfigFormat
	RenderWithConfigSections bool
}

// PipelineParams request payload for querying a single pipeline.
type PipelineParams struct {
	ConfigFormat             *ConfigFormat
	RenderWithConfigSections bool
}

// UpdatedPipeline response payload after updating a pipeline successfully.
type UpdatedPipeline struct {
	AddedPorts   []PipelinePort `json:"addedPorts"`
	RemovedPorts []PipelinePort `json:"removedPorts"`

	// Pipeline checks that have been created/updated based on AutoCreatePreChecksFromConfig changes.
	AddedChecks   []PipelineCheck `json:"addedChecks"`
	RemovedChecks []PipelineCheck `json:"removedChecks"`
}

// UpdatePipelineClusterObjects update cluster objects associated to a pipeline.
type UpdatePipelineClusterObjects struct {
	ClusterObjectsIDs []string `json:"clusterObjectsIDs"`
}

// PipelineClusterObjectsParams request payload to filter cluster objects belonging to a pipeline.
type PipelineClusterObjectsParams struct {
	Last   *uint
	Before *string
}
