package types

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleManager Role = "manager"
	RoleTeam    Role = "team"
	RoleViewer  Role = "viewer"
)

func (r Role) OK() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleTeam, RoleViewer:
		return true
	}
	return false
}
