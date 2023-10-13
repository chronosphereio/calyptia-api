package types

import "time"

// FleetFile model.
type FleetFile struct {
	ID        string    `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name"`
	Contents  []byte    `json:"contents" yaml:"contents"`
	CreatedAt time.Time `json:"createdAt" yaml:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" yaml:"updatedAt"`
}

// CreateFleetFile request payload for creating a new fleet file.
type CreateFleetFile struct {
	Name     string `json:"name"`
	Contents []byte `json:"contents"`
}

// FleetFilesParams request payload for querying the fleet files.
type FleetFilesParams struct {
	Last   *uint
	Before *string
}

// FleetFiles paginated list.
type FleetFiles struct {
	Items     []FleetFile
	EndCursor *string
}

// UpdateFleetFile request payload for updating a fleet file.
type UpdateFleetFile struct {
	Name     *string `json:"name"`
	Contents *[]byte `json:"contents"`
}
