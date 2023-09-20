package errs

const (
	InvalidPipelineSecretID      = InvalidArgumentError("invalid pipeline secret ID")
	PipelineSecretAlreadyExists  = ConflictError("pipeline secret already exists")
	PipelineSecretNotFound       = NotFoundError("pipeline secret not found")
	PipelineSecretsQuotaExceeded = PermissionDeniedError("pipeline secrets quota exceeded")
)

const (
	InvalidSecretKey    = InvalidArgumentError("invalid secret key")
	SecretSizeExceeded  = InvalidArgumentError("secret size exceeded")
	SecretValueRequired = InvalidArgumentError("secret value required")
)

func NewUndefinedPipelineSecretError(s string) *DetailedError {
	return &DetailedError{Msg: "undefined pipeline secret", Detail: &s, Parent: NotFound}
}
