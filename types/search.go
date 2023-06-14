package types

type SearchQuery struct {
	ProjectID string
	Resource  SearchResource
	Term      string
}

type SearchResource string

const (
	SearchResourceCoreInstance    SearchResource = "core_instance"
	SearchResourceResourceProfile SearchResource = "resource_profile"
	SearchResourcePipeline        SearchResource = "pipeline"
	SearchResourceAgent           SearchResource = "agent"
	SearchResourceClusterObject   SearchResource = "cluster_object"
)

type SearchResult struct {
	ID   string
	Name string
}
