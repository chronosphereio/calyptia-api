package types

import "time"

// ConfigSectionKind enum with known fluent-bit section types.
type ConfigSectionKind string

const (
	SectionKindService ConfigSectionKind = "service"
	SectionKindInput   ConfigSectionKind = "input"
	SectionKindFilter  ConfigSectionKind = "filter"
	SectionKindOutput  ConfigSectionKind = "output"
)

func (k ConfigSectionKind) OK() bool {
	for _, allowed := range [...]ConfigSectionKind{
		SectionKindService,
		SectionKindInput,
		SectionKindFilter,
		SectionKindOutput,
	} {
		if k == allowed {
			return true
		}
	}
	return false
}

// ConfigSection model represents a fluent-bit config section that can be reused
// across pipelines on a project.
type ConfigSection struct {
	ID               string            `json:"id" yaml:"id"`
	ProjectID        string            `json:"projectID" yaml:"projectID"`
	Kind             ConfigSectionKind `json:"kind" yaml:"kind"`
	Properties       Pairs             `json:"properties" yaml:"properties"`
	ProcessingRuleID *string           `json:"processingRuleID" yaml:"processingRuleID"`
	CreatedAt        time.Time         `json:"createdAt" yaml:"createdAt"`
	UpdatedAt        time.Time         `json:"updatedAt" yaml:"updatedAt"`
}

func (s ConfigSection) Name() string {
	name, _ := s.Properties.Get("Name")
	if str, ok := name.(string); ok {
		return str
	}
	return ""
}

// CreateConfigSection request payload for creating a new
// fluent-bit config section on a project.
type CreateConfigSection struct {
	Kind                       ConfigSectionKind `json:"kind"`
	Properties                 Pairs             `json:"properties"`
	ProcessingRuleID           *string           `json:"-"`
	PipelineClusterObjectSetID *string           `json:"-"`
	SkipConfigValidation       bool              `json:"skipConfigValidation"`
}

// ConfigSectionsParams request payload for querying
// the fluent-bit config sections.
type ConfigSectionsParams struct {
	Last                   *uint
	Before                 *string
	IncludeProcessingRules bool
}

// ConfigSections paginated list.
type ConfigSections struct {
	Items     []ConfigSection `json:"items" yaml:"items"`
	EndCursor *string         `json:"endCursor" yaml:"endCursor"`
}

// UpdateConfigSection request payload for updating a
// fluent-bit config section on a project.
type UpdateConfigSection struct {
	Properties           *Pairs `json:"properties"`
	SkipConfigValidation bool   `json:"skipConfigValidation"`
}
