package types

import (
	"time"
)

type ClusterObjectKind string

var (
	ClusterObjectKindNamespace ClusterObjectKind = "namespace"
	AllValidClusterObjectKinds                   = [...]ClusterObjectKind{
		ClusterObjectKindNamespace,
	}
)

// ValidClusterObjectKind checks if the given status is a valid cluster object type.
func ValidClusterObjectKind(t ClusterObjectKind) bool {
	for _, ClusterObjectKind := range AllValidClusterObjectKinds {
		if t == ClusterObjectKind {
			return true
		}
	}
	return false
}

// ClusterObject type for core_instance cluster objects.
type ClusterObject struct {
	ID        string            `json:"id" yaml:"id" db:"id"`
	Name      string            `json:"name" yaml:"name" db:"name"`
	Kind      ClusterObjectKind `json:"kind" yaml:"kind" db:"kind"`
	CreatedAt time.Time         `json:"createdAt" yaml:"createdAt" db:"created_at"`
	UpdatedAt time.Time         `json:"updatedAt" yaml:"updatedAt" db:"updated_at"`
}

// ClusterObjects paginated list.
type ClusterObjects struct {
	Items     []ClusterObject `json:"items" yaml:"items"`
	EndCursor *string         `json:"endCursor" yaml:"endCursor"`
}

// ClusterObjectParams request payload for querying the core_instance cluster objects.
type ClusterObjectParams struct {
	Last   *uint
	Name   *string
	Kind   *ClusterObjectKind
	Before *string
}

// UpdateClusterObject request payload for updating a core_instance cluster object.
type UpdateClusterObject struct {
	Name *string            `json:"name" yaml:"name"`
	Kind *ClusterObjectKind `json:"kind" yaml:"kind"`
}

// CreateClusterObject request payload for creating a core_instance cluster object.
type CreateClusterObject struct {
	Name string            `json:"name" yaml:"name"`
	Kind ClusterObjectKind `json:"kind" yaml:"kind"`
}

// CreatedClusterObject response payload after creating a core_instance cluster object.
type CreatedClusterObject struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
}
