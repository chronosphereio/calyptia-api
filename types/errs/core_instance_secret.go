package errs

const (
	CoreInstanceSecretAlreadyExists = ConflictError("core instance secret already exists")
	CoreInstanceSecretNotFound      = NotFoundError("core instance secret not found")
	InvalidCoreInstanceSecretID     = InvalidArgumentError("invalid core instance secret ID")
)
