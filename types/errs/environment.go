package errs

const (
	EnvironmentNotFound    = NotFoundError("environment not found")
	InvalidEnvironmentName = InvalidArgumentError("invalid environment name")
	EnvironmentNameTaken   = ConflictError("environment name taken")
	InvalidEnvironmentID   = InvalidArgumentError("invalid environment ID")
)
