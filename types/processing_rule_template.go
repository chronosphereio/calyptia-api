package types

import "time"

type ProcessingRuleTemplate struct {
	ID         string            `json:"id" yaml:"id" db:"id"`
	ProjectID  string            `json:"projectID" yaml:"projectID" db:"project_id"`
	Name       string            `json:"name" yaml:"name" db:"name"`
	Definition ProcessingRuleDef `json:"definition" yaml:"definition" db:"definition"`
	CreatedAt  time.Time         `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt  time.Time         `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

type ProcessingRuleDef struct {
	Name  string `json:"name" yaml:"name"`
	Match string `json:"match" yaml:"match"`
	Code  string `json:"code" yaml:"code"`
	Call  string `json:"call" yaml:"call"`
}

type StaticProcessingRuleTemplate struct {
	Name       string            `json:"name" yaml:"name"`
	Definition ProcessingRuleDef `json:"definition" yaml:"definition"`
}

type CreateProcessingRuleTemplate struct {
	ProjectID  string            `json:"-"`
	Name       string            `json:"name"`
	Definition ProcessingRuleDef `json:"definition"`
}

type ListProcessingRuleTemplates struct {
	ProjectID string
	Last      *uint
	Before    *string
	Name      *string
}

type ProcessingRuleTemplates struct {
	Items     []ProcessingRuleTemplate       `json:"items" yaml:"items"`
	EndCursor *string                        `json:"endCursor" yaml:"endCursor"`
	Count     uint                           `json:"count" yaml:"count"`
	Static    []StaticProcessingRuleTemplate `json:"static" yaml:"static"`
}

type UpdateProcessingRuleTemplate struct {
	TemplateID string             `json:"-"`
	Name       *string            `json:"name"`
	Definition *ProcessingRuleDef `json:"definition"`
}
