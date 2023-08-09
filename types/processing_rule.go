package types

import "time"

// ProcessingRule defines a set of actions that
// eventually will get translated into a single fluent-bif filter.
// This filter is stored as a config section associated with a pipeline.
type ProcessingRule struct {
	ID              string                 `json:"id" yaml:"id"`
	PipelineID      string                 `json:"pipelineID" yaml:"pipelineID"`
	ConfigSectionID string                 `json:"configSectionID" yaml:"configSectionID"`
	FileID          string                 `json:"fileID" yaml:"fileID"`
	Match           string                 `json:"match" yaml:"match"`
	IsMatchRegexp   bool                   `json:"isMatchRegexp" yaml:"isMatchRegexp"`
	Language        ProcessingRuleLanguage `json:"language" yaml:"language"`
	Actions         []RuleAction           `json:"actions" yaml:"actions"`
	CreatedAt       time.Time              `json:"createdAt" yaml:"createdAt"`
	UpdatedAt       time.Time              `json:"updatedAt" yaml:"updatedAt"`
}

// ProcessingRuleLanguage enum of scripting languages
// a processing rule supports.
type ProcessingRuleLanguage string

// ProcessingRuleLanguageLua will produce a filter in Lua.
const ProcessingRuleLanguageLua ProcessingRuleLanguage = "lua"

// RuleAction within a processing rule.
// Each action is processed in order.
type RuleAction struct {
	Kind        RuleActionKind `json:"kind" yaml:"kind"`
	Description string         `json:"description" yaml:"description"`
	Enabled     bool           `json:"enabled" yaml:"enabled"`

	Selectors []LogSelector `json:"selectors" yaml:"selectors"`

	// oneof the following depending on Kind:

	Add      *LogAttr `json:"add,omitempty" yaml:"add,omitempty"`
	RenameTo *string  `json:"renameTo,omitempty" yaml:"renameTo,omitempty"`
	CopyAs   *string  `json:"copyAs,omitempty" yaml:"copyAs,omitempty"`
	MaskWith *string  `json:"maskWith,omitempty" yaml:"maskWith,omitempty"`
}

// RuleActionKind enum of the different action kinds a processing rule can have.
type RuleActionKind string

const (
	// RuleActionKindAdd adds a key-value pair to the log record.
	RuleActionKindAdd RuleActionKind = "add"
	// RuleActionKindRename renames the matching key into the new key.
	RuleActionKindRename RuleActionKind = "rename"
	// RuleActionKindCopy copies the matching key into the new key.
	// Conserving both.
	RuleActionKindCopy RuleActionKind = "copy"
	// RuleActionKindMask causes the value associated
	// with the matching key to be replaced with `redacted`.
	RuleActionKindMask RuleActionKind = "mask"
	// RuleActionKindRemove removes a key-value pair
	// from the log record using its key.
	RuleActionKindRemove RuleActionKind = "remove"
	// RuleActionKindSkip causes a log record to be skipped entirely
	// using its key.
	RuleActionKindSkip RuleActionKind = "skip"
)

// LogSelector used to match a log entry.
// Example:
//
//   - Source.kind=key Op=equal Target=foo
//     matches a log with a key equal to "foo" on it.
//   - Source.kind=value Op=equal Target=bar
//     matches a log with a value equal to "bar" on it.
type LogSelector struct {
	Kind LogSelectorKind   `json:"kind" yaml:"kind"`
	Op   LogSelectorOpKind `json:"op" yaml:"op"`
	Expr string            `json:"expr" yaml:"expr"`
}

// LogSelectorKind enum.
type LogSelectorKind string

const (
	// LogSelectorKindKey matches a log using some of its keys.
	LogSelectorKindKey LogSelectorKind = "key"
)

// LogSelectorOpKind enum of the supported operations a selector.
type LogSelectorOpKind string

const (
	// LogSelectorOpKindEqual matches a log key/value equally.
	LogSelectorOpKindEqual LogSelectorOpKind = "equal"
)

// LogAttr its the key-value pair in a log record.
type LogAttr struct {
	Key   string `json:"key" yaml:"key"`
	Value string `json:"value" yaml:"value"`
}

// CreateProcessingRule request payload.
type CreateProcessingRule struct {
	PipelineID    string                 `json:"-"`
	Match         string                 `json:"match"`
	IsMatchRegexp bool                   `json:"isMatchRegexp"`
	Language      ProcessingRuleLanguage `json:"language"`
	Actions       []RuleAction           `json:"actions"`
}

// CreatedProcessingRule response payload.
type CreatedProcessingRule struct {
	ID              string    `json:"id" yaml:"id"`
	ConfigSectionID string    `json:"configSectionID" yaml:"configSectionID"`
	FileID          string    `json:"fileID" yaml:"fileID"`
	CreatedAt       time.Time `json:"createdAt" yaml:"createdAt"`
}

// ProcessingRulesParams request payload for querying processing rules.
type ProcessingRulesParams struct {
	PipelineID string
	Last       *uint
	Before     *string
}

// ProcessingRules paginated list.
type ProcessingRules struct {
	Items     []ProcessingRule `json:"items" yaml:"items"`
	EndCursor *string          `json:"endCursor" yaml:"endCursor"`
}

// UpdateProcessingRule request payload.
type UpdateProcessingRule struct {
	ProcessingRuleID string                  `json:"-"`
	Match            *string                 `json:"match,omitempty"`
	IsMatchRegexp    *bool                   `json:"isMatchRegexp,omitempty"`
	Language         *ProcessingRuleLanguage `json:"language,omitempty"`
	Actions          *[]RuleAction           `json:"actions,omitempty"`
}

// PreviewProcessingRule request payload to run and preview the input/output of
// a processing rule.
// Given some sample logs, and some actions, it shows the output log records.
type PreviewProcessingRule struct {
	Language ProcessingRuleLanguage `json:"language"`
	Actions  []RuleAction           `json:"actions"`
	Logs     []FluentBitLog         `json:"logs"`
}
