package types

import "time"

// ResourceProfile model.
type ResourceProfile struct {
	ID                     string    `json:"id" yaml:"id"`
	Name                   string    `json:"name" yaml:"name"`
	StorageMaxChunksUp     int       `json:"storageMaxChunksUp" yaml:"storageMaxChunksUp"`
	StorageSyncFull        bool      `json:"storageSyncFull" yaml:"storageSyncFull"`
	StorageBacklogMemLimit string    `json:"storageBacklogMemLimit" yaml:"storageBacklogMemLimit"`
	StorageVolumeSize      string    `json:"storageVolumeSize" yaml:"storageVolumeSize"`
	StorageMaxChunksPause  bool      `json:"storageMaxChunksPause" yaml:"storageMaxChunksPause"`
	CPUBufferWorkers       int       `json:"cpuBufferWorkers" yaml:"cpuBufferWorkers"`
	CPULimit               string    `json:"cpuLimit" yaml:"cpuLimit"`
	CPURequest             string    `json:"cpuRequest" yaml:"cpuRequest"`
	MemoryLimit            string    `json:"memoryLimit" yaml:"memoryLimit"`
	MemoryRequest          string    `json:"memoryRequest" yaml:"memoryRequest"`
	CreatedAt              time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt              time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// ResourceProfiles paginated list.
type ResourceProfiles struct {
	Items     []ResourceProfile
	EndCursor *string
}

const (
	// ResourceProfileHighPerformanceGuaranteedDelivery is one of the 3 default resource profiles created with an aggregator.
	ResourceProfileHighPerformanceGuaranteedDelivery = "high-performance-guaranteed-delivery"
	// ResourceProfileHighPerformanceOptimalThroughput is one of the 3 default resource profiles created with an aggregator.
	ResourceProfileHighPerformanceOptimalThroughput = "high-performance-optimal-throughput"
	// ResourceProfileBestEffortLowResource is one of the 3 default resource profiles created with an aggregator. This is the default one.
	ResourceProfileBestEffortLowResource = "best-effort-low-resource"
)

// DefaultResourceProfileName is the default resource profile used when creating pipelines.
const DefaultResourceProfileName = ResourceProfileBestEffortLowResource

// CreateResourceProfile request payload for creating a resource profile.
type CreateResourceProfile struct {
	Name                   string `json:"name"`
	StorageMaxChunksUp     uint   `json:"storageMaxChunksUp"`
	StorageSyncFull        bool   `json:"storageSyncFull"`
	StorageBacklogMemLimit string `json:"storageBacklogMemLimit"`
	StorageVolumeSize      string `json:"storageVolumeSize"`
	StorageMaxChunksPause  bool   `json:"storageMaxChunksPause"`
	CPUBufferWorkers       uint   `json:"cpuBufferWorkers"`
	CPULimit               string `json:"cpuLimit"`
	CPURequest             string `json:"cpuRequest"`
	MemoryLimit            string `json:"memoryLimit"`
	MemoryRequest          string `json:"memoryRequest"`
}

// ResourceProfilesParams request payload for querying resource profiles.
type ResourceProfilesParams struct {
	Last   *uint
	Before *string
}

// UpdateResourceProfile request payload for updating a resource profile.
type UpdateResourceProfile struct {
	Name                   *string `json:"name"`
	StorageMaxChunksUp     *uint   `json:"storageMaxChunksUp"`
	StorageSyncFull        *bool   `json:"storageSyncFull"`
	StorageMaxChunksPause  *bool   `json:"storageMaxChunksPause"`
	StorageBacklogMemLimit *string `json:"storageBacklogMemLimit"`
	StorageVolumeSize      *string `json:"storageVolumeSize"`
	CPUBufferWorkers       *uint   `json:"cpuBufferWorkers"`
	CPULimit               *string `json:"cpuLimit"`
	CPURequest             *string `json:"cpuRequest"`
	MemoryLimit            *string `json:"memoryLimit"`
	MemoryRequest          *string `json:"memoryRequest"`
}
