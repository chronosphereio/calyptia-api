package types

import "time"

// FluentBitPluginKind enum with known fluent-bit plugin types.
type FluentBitPluginKind string

const (
	FluentBitPluginKindInput  FluentBitPluginKind = "input"
	FluentBitPluginKindFilter FluentBitPluginKind = "filter"
	FluentBitPluginKindOutput FluentBitPluginKind = "output"
)

// FluentBitPluginConfig model represents a fluent-bit plugin that can be reused
// across pipelines on a project.
type FluentBitPluginConfig struct {
	ID         string                 `json:"id" yaml:"id"`
	ProjectID  string                 `json:"projectID" yaml:"projectID"`
	Kind       FluentBitPluginKind    `json:"kind" yaml:"kind"`
	Name       string                 `json:"name" yaml:"name"`
	Properties map[string]interface{} `json:"properties" yaml:"properties"`
	CreatedAt  time.Time              `json:"createdAt" yaml:"createdAt"`
	UpdatedAt  time.Time              `json:"updatedAt" yaml:"updatedAt"`
}

// CreateFluentBitPluginConfig request payload for creating a new
// fluent-bit plugin config on a project.
type CreateFluentBitPluginConfig struct {
	Kind       FluentBitPluginKind    `json:"kind"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
}

// CreatedFluentBitPluginConfig response payload after creating
// a fluent-bit plugin config successfully.
type CreatedFluentBitPluginConfig struct {
	ID        string    `json:"id" yaml:"id"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
}

// FluentBitPluginConfigsParams request payload for querying
// the fluent-bit plugin configs.
type FluentBitPluginConfigsParams struct {
	Last   *uint
	Before *string
}

// FluentBitPluginConfigs paginated list.
type FluentBitPluginConfigs struct {
	Items     []FluentBitPluginConfig `json:"items" yaml:"items"`
	EndCursor *string                 `json:"endCursor" yaml:"endCursor"`
}
