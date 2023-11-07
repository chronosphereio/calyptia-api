package errs

const (
	PipelineLogNotFound  = NotFoundError("pipeline log not found")
	InvalidPipelineLogs  = InvalidArgumentError("invalid pipeline logs")
	InvalidPipelineLogID = InvalidArgumentError("invalid pipeline log ID")
)
