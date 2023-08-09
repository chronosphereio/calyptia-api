package types

import (
	"time"
)

type Fleet struct {
	ID                  string           `json:"id" yaml:"id"`
	ProjectID           string           `json:"projectID" yaml:"projectID"`
	Name                string           `json:"name" yaml:"name"`
	MinFluentBitVersion string           `json:"minFluentBitVersion" yaml:"minFluentBitVersion"`
	RawConfig           string           `json:"rawConfig" yaml:"rawConfig"`
	ConfigFormat        ConfigFormat     `json:"configFormat" yaml:"configFormat"`
	Tags                []string         `json:"tags" yaml:"tags"`
	AgentsCount         FleetAgentsCount `json:"agentsCount" yaml:"agentsCount"`
	CreatedAt           time.Time        `json:"createdAt" yaml:"createdAt"`
	UpdatedAt           time.Time        `json:"updatedAt" yaml:"updatedAt"`
}

type FleetAgentsCount struct {
	Active     uint `json:"active" yaml:"active"`
	Inactive   uint `json:"inactive" yaml:"inactive"`
	WithErrors uint `json:"errors" yaml:"withErrors"`
}

type CreateFleet struct {
	ProjectID           string       `json:"-"`
	Name                string       `json:"name"`
	MinFluentBitVersion string       `json:"minFluentBitVersion"`
	RawConfig           string       `json:"rawConfig"`
	ConfigFormat        ConfigFormat `json:"configFormat"`
	Tags                []string     `json:"tags"`

	SkipConfigValidation bool `json:"skipConfigValidation"`
}

type FleetsParams struct {
	ProjectID    string
	Name         *string
	TagsQuery    *string
	Last         *uint
	Before       *string
	ConfigFormat *ConfigFormat

	tags *[]string
}

// FleetParams request payload for querying a single fleet.
type FleetParams struct {
	FleetID      string
	ConfigFormat *ConfigFormat
}

func (p FleetsParams) Tags() ([]string, bool) {
	if p.tags == nil {
		return nil, false
	}

	return *p.tags, true
}

func (p *FleetsParams) SetTags(tt []string) {
	p.tags = &tt
}

type Fleets struct {
	Items     []Fleet `json:"items" yaml:"items"`
	EndCursor *string `json:"endCursor" yaml:"endCursor"`
}

type UpdateFleet struct {
	ID           string        `json:"-"`
	Name         *string       `json:"name"`
	RawConfig    *string       `json:"rawConfig"`
	ConfigFormat *ConfigFormat `json:"configFormat"`
	Tags         *[]string     `json:"tags"`

	SkipConfigValidation bool `json:"skipConfigValidation"`
}

type FleetConfigParams struct {
	ConfigFormat *ConfigFormat
}
