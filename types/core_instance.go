package types

import (
	"bytes"
	"encoding/json"
	"time"
)

const (
	DefaultCoreInstanceVersion = "v0.1.12"
)

// CoreInstance model.
type CoreInstance struct {
	ID                    string               `json:"id" yaml:"id"`
	Token                 string               `json:"token" yaml:"token"`
	Name                  string               `json:"name" yaml:"name"`
	ClusterLoggingEnabled bool                 `json:"clusterLoggingEnabled" yaml:"clusterLoggingEnabled"`
	EnvironmentName       string               `json:"environmentName" yaml:"environmentName"`
	Version               string               `json:"version" yaml:"version"`
	PipelinesCount        uint                 `json:"pipelinesCount" yaml:"pipelinesCount"`
	Tags                  []string             `json:"tags" yaml:"tags"`
	Metadata              CoreInstanceMetadata `json:"metadata" yaml:"metadata"`
	Status                CoreInstanceStatus   `json:"status" yaml:"status"`
	SkipServiceCreation   bool                 `json:"skipServiceCreation" yaml:"skipServiceCreation"`
	//  Pointer to a string (*string) allows the value to be nullable, which means it can be assigned a NULL value from the database.
	//  If a value is not nullable, such as a regular string, it cannot be assigned a NULL value from the database and would result in the error "cannot scan NULL into string."
	Image     *string   `json:"image" yaml:"image"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

type CoreInstanceStatus string

const (
	CoreInstanceStatusWaiting     CoreInstanceStatus = "waiting"
	CoreInstanceStatusRunning     CoreInstanceStatus = "running"
	CoreInstanceStatusUnreachable CoreInstanceStatus = "unreachable"
)

// CoreInstance ping constants.
const (
	// CoreInstanceNextPing is the time between pings to a core instance.
	CoreInstanceNextPing = time.Second * 30
	// CoreInstanceNextPingDelta is the extra time acceptable for a ping to be delayed.
	CoreInstanceNextPingDelta = time.Second * 5
	// CoreInstanceNextPingTimeout is the time after a core instance is considered "unreachable".
	CoreInstanceNextPingTimeout = CoreInstanceNextPing + CoreInstanceNextPingDelta
)

// CoreInstances paginated list.
type CoreInstances struct {
	Items     []CoreInstance `json:"items"`
	EndCursor *string        `json:"endCursor"`
	Count     int            `json:"count"`
}

// CoreInstancePingResponse response from a core instance ping call.
type CoreInstancePingResponse struct {
	NextPing time.Duration `json:"nextPing"`
}

// CoreInstanceMetadata See: https://github.com/fluent/fluent-bit/blob/d059a5a5dca6aff4ff5d0694887355480d6f2c1d/plugins/out_calyptia/calyptia.c#L206-L229
// Those are the only supported metadata fields that will be marshaled by the Calyptia Cloud API, please send a PR if further
// fields are required.
// This represents a blob of JSON that will be stored as it is in the Cloud database.
// That's why all fields have an omitempty json tag; to avoid filling the blob with empty data.
type CoreInstanceMetadata struct {
	// Notice that all of these are embedded on purpose since
	// metadata is flattened.
	MetadataK8S
	MetadataAWS
	MetadataGCP
}

// UnmarshalJSON deserializes JSON data into a CoreInstanceMetadata struct.
// Implements the json.Unmarshaler interface.
//
// If the input data is nil, UnmarshalJSON returns nil without modifying the receiver.
//
// If the input data is "null", an empty CoreInstanceMetadata struct is assigned to the receiver.
//
// If the input data is an empty JSON object "{}", an empty CoreInstanceMetadata struct is assigned to the receiver.
//
// If the input data is any other valid JSON object, it is first unmarshalled into a copy of CoreInstanceMetadata using the
// standard library's json.Unmarshal function, and then the copy is assigned to the receiver.
func (m *CoreInstanceMetadata) UnmarshalJSON(data []byte) error {
	if len(data) == 0 || bytes.Equal(data, []byte("null")) || bytes.Equal(data, []byte("{}")) {
		*m = CoreInstanceMetadata{}
		return nil
	}

	// This is a workaround to avoid a recursive call to UnmarshalJSON.
	type Copy CoreInstanceMetadata
	var dest Copy
	if err := json.Unmarshal(data, &dest); err != nil {
		return err
	}

	*m = CoreInstanceMetadata(dest)
	return nil
}

// MetadataK8S See: https://github.com/kubernetes/website/blob/60390ff3c0ef0043a58568ad2e4c2b7634028074/content/en/examples/pods/inject/dapi-volume.yaml#L5
// For further cluster information data that can be included check: https://pkg.go.dev/k8s.io/client-go/discovery#DiscoveryClient.
type MetadataK8S struct {
	Namespace       string `json:"k8s.namespace,omitempty"`
	ClusterName     string `json:"k8s.cluster_name,omitempty"`
	Zone            string `json:"k8s.zone,omitempty"`
	ClusterVersion  string `json:"k8s.cluster_version,omitempty"`
	ClusterPlatform string `json:"k8s.cluster_platform,omitempty"`
}

// MetadataAWS See: https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-retrieval.html
type MetadataAWS struct {
	AMIID           string `json:"aws.ami_id,omitempty"`
	AccountID       string `json:"aws.account_id,omitempty"`
	Hostname        string `json:"aws.hostname,omitempty"`
	VPCID           string `json:"aws.vpc_id,omitempty"`
	PrivateIPv4     string `json:"aws.private_ipv4,omitempty"`
	PublicIPv4      string `json:"aws.public_ipv4,omitempty"`
	EC2InstanceID   string `json:"aws.ec2_instance_id,omitempty"`
	EC2InstanceType string `json:"aws.ec2_instance_type,omitempty"`
	AZ              string `json:"aws.az,omitempty"`
}

// MetadataGCP See: https://cloud.google.com/compute/docs/metadata/default-metadata-values
type MetadataGCP struct {
	ProjectNumber uint64 `json:"gcp.project_number,omitempty"`
	ProjectID     string `json:"gcp.project_id,omitempty"`
	InstanceID    uint64 `json:"gcp.instance_id,omitempty"`
	InstanceImage string `json:"gcp.instance_image,omitempty"`
	MachineType   string `json:"gcp.machine_type,omitempty"`
	InstanceName  string `json:"gcp.instance_name,omitempty"`
	Zone          string `json:"gcp.zone,omitempty"`
}

// CreateCoreInstance request payload for creating a new core instance.
type CreateCoreInstance struct {
	Name                    string               `json:"name"`
	Version                 string               `json:"version"`
	AddHealthCheckPipeline  bool                 `json:"addHealthCheckPipeline"`
	HealthCheckPipelinePort uint                 `json:"healthCheckPipelinePort"`
	ClusterLogging          bool                 `json:"clusterLogging"`
	Tags                    []string             `json:"tags"`
	Image                   *string              `json:"image"`
	Metadata                CoreInstanceMetadata `json:"metadata"`
	EnvironmentID           string               `json:"environmentID"`
	SkipServiceCreation     bool                 `json:"skipServiceCreation"`
}

// CreatedCoreInstance response payload after creating a core instance successfully.
type CreatedCoreInstance struct {
	ID                     string            `json:"id"`
	Token                  string            `json:"token"`
	PrivateRSAKey          []byte            `json:"privateRSAKey"`
	PublicRSAKey           []byte            `json:"publicRSAKey"`
	Name                   string            `json:"name"`
	Version                string            `json:"version"`
	Image                  string            `json:"image"`
	CreatedAt              time.Time         `json:"createdAt"`
	Tags                   []string          `json:"tags"`
	HealthCheckPipeline    *Pipeline         `json:"healthCheckPipeline"`
	ClusterLoggingPipeline *Pipeline         `json:"clusterLoggingPipeline"`
	ResourceProfiles       []ResourceProfile `json:"resourceProfiles"`
	EnvironmentName        string            `json:"environmentName"`
	SkipServiceCreation    bool              `json:"skipServiceCreation"`
}

// CoreInstancesParams request payload for querying core instances.
type CoreInstancesParams struct {
	Last          *uint
	Before        *string
	Name          *string
	Tags          *string
	EnvironmentID *string
}

// UpdateCoreInstance request payload for updating a core instance.
type UpdateCoreInstance struct {
	Name                *string               `json:"name"`
	Version             *string               `json:"version"`
	EnvironmentID       *string               `json:"environmentID"`
	ClusterLogging      *bool                 `json:"clusterLogging"`
	Tags                *[]string             `json:"tags"`
	Metadata            *CoreInstanceMetadata `json:"metadata"`
	SkipServiceCreation *bool                 `json:"skipServiceCreation"`
	Image               *string               `json:"image"`
}
