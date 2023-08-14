package errs

const (
	InvalidPermission   = InvalidArgumentError("invalid permission")
	DuplicatePermission = InvalidArgumentError("duplicate permission")
)
