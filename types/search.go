package types

type SearchQuery struct {
	ProjectID string
	Resource  SearchResource
	Term      string
}

type SearchResource string

const (
	SearchResourceAgent           SearchResource = "agent"
	SearchResourceClusterObject   SearchResource = "cluster_object"
	SearchResourceConfigSection   SearchResource = "config_section"
	SearchResourceCoreInstance    SearchResource = "core_instance"
	SearchResourceEnvironment     SearchResource = "environment"
	SearchResourceFleet           SearchResource = "fleet"
	SearchResourceMember          SearchResource = "member"
	SearchResourcePipeline        SearchResource = "pipeline"
	SearchResourcePipelineSecret  SearchResource = "pipeline_secret"
	SearchResourceResourceProfile SearchResource = "resource_profile"
	SearchResourceTraceSession    SearchResource = "trace_session"
)

type SearchResult struct {
	ID   string `json:"id" yaml:"id" db:"id"`
	Name string `json:"name" yaml:"name" db:"name"`
}
