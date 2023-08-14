package errs

const (
	InvalidResourceProfileID     = InvalidArgumentError("invalid resource profile ID")
	InvalidResourceProfileName   = InvalidArgumentError("invalid resource profile name")
	ResourceProfileStillInUse    = PermissionDeniedError("resource profile still in use by pipelines")
	ResourceProfileAlreadyExists = ConflictError("resource profile name already exists")
	ResourceProfileNotFound      = NotFoundError("resource profile not found")
)
