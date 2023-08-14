package errs

const (
	InvalidPipelineConfig       = InvalidArgumentError("invalid pipeline config")
	InvalidPipelineConfigFormat = InvalidArgumentError("invalid pipeline config format")
	PipelineConfigNotFound      = NotFoundError("pipeline config not found")
)
