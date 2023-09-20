package errs

const (
	CoreInstanceFileAlreadyExists = ConflictError("core instance file already exists")
	CoreInstanceFileNotFound      = NotFoundError("core instance file not found")
	InvalidCoreInstanceFileID     = InvalidArgumentError("invalid core instance file ID")
)
