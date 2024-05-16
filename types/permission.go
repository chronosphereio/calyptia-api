package types

const (
	PermCreateAll = "create:*"
	PermReadAll   = "read:*"
	PermUpdateAll = "update:*"
	PermDeleteAll = "delete:*"
)

func AllPermissions() []string {
	return []string{
		PermCreateAll,
		PermReadAll,
		PermUpdateAll,
		PermDeleteAll,
	}
}
