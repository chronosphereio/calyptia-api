package types

import "time"

const ErrMemberAlreadyExists = ConflictError("member already exists")

// Membership model.
type Membership struct {
	ID        string           `json:"id" yaml:"id"`
	Roles     []MembershipRole `json:"roles" yaml:"roles"`
	CreatedAt time.Time        `json:"createdAt" yaml:"createdAt"`

	User *User `json:"user" yaml:"user"`
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
	Last *uint64
}
