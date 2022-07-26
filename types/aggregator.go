package types

import (
	"encoding/json"
	"time"
)

const (
	DefaultAggregatorVersion = "v0.1.12"
)

// Aggregator model.
type Aggregator struct {
	ID              string           `json:"id" yaml:"id"`
	Token           string           `json:"token" yaml:"token"`
	Name            string           `json:"name" yaml:"name"`
	EnvironmentName string           `json:"environmentName" yaml:"environmentName"`
	Version         string           `json:"version" yaml:"version"`
	PipelinesCount  uint64           `json:"pipelinesCount" yaml:"pipelinesCount"`
	Tags            []string         `json:"tags" yaml:"tags"`
	Metadata        *json.RawMessage `json:"metadata" yaml:"metadata"`
	Status          AggregatorStatus `json:"status" yaml:"status"`
	CreatedAt       time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

type AggregatorStatus string

const (
	AggregatorStatusWaiting     AggregatorStatus = "waiting"
	AggregatorStatusRunning     AggregatorStatus = "running"
	AggregatorStatusUnreachable AggregatorStatus = "unreachable"
)

// Aggregator ping constants.
const (
	// AggregatorNextPing is the time between pings to an aggregator.
	AggregatorNextPing = time.Second * 30
	// AggregatorNextPingDelta is the extra time acceptable for a ping to be delayed.
	AggregatorNextPingDelta = time.Second * 5
	// AggregatorNextPingTimeout is the time after an aggregator is considered "unreachable".
	AggregatorNextPingTimeout = AggregatorNextPing + AggregatorNextPingDelta
)

// Aggregators paginated list.
type Aggregators struct {
	Items     []Aggregator
	EndCursor *string
}

// AggregatorPingResponse response from an aggregator ping call.
type AggregatorPingResponse struct {
	NextPing time.Duration `json:"nextPing"`
}

// AggregatorMetadata See: https://github.com/fluent/fluent-bit/blob/d059a5a5dca6aff4ff5d0694887355480d6f2c1d/plugins/out_calyptia/calyptia.c#L206-L229
// Those are the only supported metadata fields that will be marshaled by the Calyptia Cloud API, please send a PR if further
// fields are required.
type AggregatorMetadata struct {
	// Notice that all of these are embedded on purpose since
	// metadata is flattened.
	MetadataK8S
	MetadataAWS
	MetadataGCP
}

// MetadataK8S See: https://github.com/kubernetes/website/blob/60390ff3c0ef0043a58568ad2e4c2b7634028074/content/en/examples/pods/inject/dapi-volume.yaml#L5
// For further cluster information data that can be included check: https://pkg.go.dev/k8s.io/client-go/discovery#DiscoveryClient.
type MetadataK8S struct {
	Namespace       string `json:"k8s.namespace"`
	ClusterName     string `json:"k8s.cluster_name"`
	Zone            string `json:"k8s.zone"`
	ClusterVersion  string `json:"k8s.cluster_version"`
	ClusterPlatform string `json:"k8s.cluster_platform"`
}

// MetadataAWS See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html
type MetadataAWS struct {
	AMIID           string `json:"aws.ami_id"`
	AccountID       string `json:"aws.account_id"`
	Hostname        string `json:"aws.hostname"`
	VPCID           string `json:"aws.vpc_id"`
	PrivateIPv4     string `json:"aws.private_ipv4"`
	PublicIPv4      string `json:"aws.public_ipv4"`
	EC2InstanceID   string `json:"aws.ec2_instance_id"`
	EC2InstanceType string `json:"aws.ec2_instance_type"`
	AZ              string `json:"aws.az"`
}

// MetadataGCP See: https://cloud.google.com/compute/docs/metadata/default-metadata-values
type MetadataGCP struct {
	ProjectNumber uint64 `json:"gcp.project_number"`
	ProjectID     string `json:"gcp.project_id"`
	InstanceID    uint64 `json:"gcp.instance_id"`
	InstanceImage string `json:"gcp.instance_image"`
	MachineType   string `json:"gcp.machine_type"`
	InstanceName  string `json:"gcp.instance_name"`
	Zone          string `json:"gcp.zone"`
}

// CreateAggregator request payload for creating a new aggregator.
type CreateAggregator struct {
	Name                    string           `json:"name"`
	Version                 string           `json:"version"`
	AddHealthCheckPipeline  bool             `json:"addHealthCheckPipeline"`
	HealthCheckPipelinePort uint             `json:"healthCheckPipelinePort"`
	Tags                    []string         `json:"tags"`
	Metadata                *json.RawMessage `json:"metadata"`
	EnvironmentID           string           `json:"environmentID"`
}

// CreatedAggregator response payload after creating an aggregator successfully.
type CreatedAggregator struct {
	ID            string    `json:"id"`
	Token         string    `json:"token"`
	PrivateRSAKey []byte    `json:"privateRSAKey"`
	PublicRSAKey  []byte    `json:"publicRSAKey"`
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	CreatedAt     time.Time `json:"createdAt"`
	Tags          []string  `json:"tags"`

	HealthCheckPipeline *Pipeline         `json:"healthCheckPipeline"`
	ResourceProfiles    []ResourceProfile `json:"resourceProfiles"`
}

// AggregatorsParams request payload for querying aggregators.
type AggregatorsParams struct {
	Last   *uint64
	Before *string
	Name   *string
	Tags   *string
}

// UpdateAggregator request payload for updating an aggregator.
type UpdateAggregator struct {
	Name          *string          `json:"name"`
	Version       *string          `json:"version"`
	EnvironmentID *string          `json:"environmentID"`
	Metadata      *json.RawMessage `json:"metadata"`
}
