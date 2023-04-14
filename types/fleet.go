package types

import (
	"time"

	fluentbitconfig "github.com/calyptia/go-fluentbit-config"
)

type Fleet struct {
	ID                  string                 `json:"id"`
	ProjectID           string                 `json:"projectID"`
	Name                string                 `json:"name"`
	MinFluentBitVersion string                 `json:"minFluentBitVersion"`
	Config              fluentbitconfig.Config `json:"config"`
	Tags                []string               `json:"tags"`
	CreatedAt           time.Time              `json:"createdAt"`
	UpdatedAt           time.Time              `json:"updatedAt"`
}

type FleetConfigParams struct {
	ConfigFormat *ConfigFormat
}

type CreateFleet struct {
	ProjectID           string                 `json:"-"`
	Name                string                 `json:"name"`
	MinFluentBitVersion string                 `json:"minFluentBitVersion"`
	Config              fluentbitconfig.Config `json:"config"`
	Tags                []string               `json:"tags"`

	SkipConfigValidation bool `json:"skipConfigValidation"`
}

type CreatedFleet struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}

type FleetsParams struct {
	ProjectID string
	Name      *string
	TagsQuery *string
	Last      *uint
	Before    *string

	tags *[]string
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
	Items     []Fleet `json:"items"`
	EndCursor *string `json:"endCursor"`
}

type UpdateFleet struct {
	ID     string                  `json:"-"`
	Name   *string                 `json:"name"`
	Config *fluentbitconfig.Config `json:"config"`
	Tags   *[]string               `json:"tags"`

	SkipConfigValidation bool `json:"skipConfigValidation"`
}

type UpdatedFleet struct {
	UpdatedAt time.Time `json:"updatedAt"`
}

type DeletedFleet struct {
	Deleted   bool       `json:"deleted"`
	DeletedAt *time.Time `json:"deletedAt"`
}
