package types

import "time"

// Membership model.
type Membership struct {
	ID        string           `json:"id" yaml:"id"`
	Roles     []MembershipRole `json:"roles" yaml:"roles"`
	CreatedAt time.Time        `json:"createdAt" yaml:"createdAt"`

	User *User `json:"user" yaml:"user"`
}

type MembershipRole string

const (
	MembershipRoleCreator MembershipRole = "creator"
	MembershipRoleAdmin   MembershipRole = "admin"
)

// MembersParams request payload for querying members.
type MembersParams struct {
	Last *uint64
}
