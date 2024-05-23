package types

import "time"

// Membership model.
type Membership struct {
	ID          string           `json:"id" yaml:"id" db:"id"`
	UserID      string           `json:"userID" yaml:"userID" db:"user_id"`
	ProjectID   string           `json:"projectID" yaml:"projectID" db:"project_id"`
	Roles       []MembershipRole `json:"roles" yaml:"roles" db:"roles"`
	Permissions []string         `json:"permissions" yaml:"permissions" db:"permissions"`
	CreatedAt   time.Time        `json:"createdAt" yaml:"createdAt" db:"created_at"`

	User *User `json:"user" yaml:"user"`
}

func (m Membership) HasRole(role MembershipRole) bool {
	return hasRole(m.Roles, role)
}

func hasRole(roles []MembershipRole, role MembershipRole) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// CanManageMembership checks if the current member has authority over the target member.
func (m Membership) CanManageMembership(target Membership) bool {
	if m.ID == target.ID || m.HasRole(MembershipRoleCreator) {
		return true
	}

	if len(m.Roles) == 0 {
		return false
	}

	if m.HasRole(MembershipRoleAdmin) {
		return !hasRole(target.Roles, MembershipRoleCreator)
	}

	return false
}

func (m Membership) HasOnlyPermission(permission string) bool {
	return len(m.Permissions) == 1 && m.Permissions[0] == permission
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

// UpdateMember request body.
type UpdateMember struct {
	MemberID    string    `json:"-"`
	Permissions *[]string `json:"permissions"`
}

type CountMembers struct {
	ProjectID string
	Role      *MembershipRole
}
