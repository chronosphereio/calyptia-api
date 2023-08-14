package errs

const (
	InvalidPipelineSecretID      = InvalidArgumentError("invalid pipeline secret ID")
	InvalidPipelineSecretKey     = InvalidArgumentError("invalid pipeline secret key")
	PipelineSecretAlreadyExists  = ConflictError("pipeline secret already exists")
	PipelineSecretNotFound       = NotFoundError("pipeline secret not found")
	PipelineSecretsQuotaExceeded = PermissionDeniedError("pipeline secrets quota exceeded")
	PipelineSecretSizeExceeded   = InvalidArgumentError("pipeline secret size exceeded")
)
