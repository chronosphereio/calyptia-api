package errs

const (
	InvalidPipelineStatus  = InvalidArgumentError("invalid pipeline status")
	PipelineStatusNotFound = NotFoundError("pipeline status not found")
)
