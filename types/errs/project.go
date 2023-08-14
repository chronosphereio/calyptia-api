package errs

const (
	InvalidProjectID   = InvalidArgumentError("invalid project ID")
	InvalidProjectName = InvalidArgumentError("invalid project name")
	ProjectNotFound    = NotFoundError("project not found")
)
