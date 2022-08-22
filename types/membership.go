package types

import "time"

// Membership model.
type Membership struct {
	ID        string           `json:"id" yaml:"id"`
	Roles     []MembershipRole `json:"roles" yaml:"roles"`
	CreatedAt time.Time        `json:"createdAt" yaml:"createdAt"`

	User *User `json:"user" yaml:"user"`
}

// Memberships paginated list.
type Memberships struct {
	Items     []Membership
	EndCursor *string
}

// MembershipRole type of membership role (creator/admin).
type MembershipRole string

const (
	// MembershipRoleCreator creator membership role.
	MembershipRoleCreator MembershipRole = "creator"
	// MembershipRoleAdmin admin membership role.
	MembershipRoleAdmin MembershipRole = "admin"
)

// MembersParams request payload for querying members.
type MembersParams struct {
	Last   *uint
	Before *string
}
