package errs

const (
	PipelineLogNotFound      = NotFoundError("pipeline log not found")
	InvalidPipelineLogs      = InvalidArgumentError("invalid pipeline logs")
	InvalidPipelineLogStatus = InvalidArgumentError("invalid pipeline log status")
	InvalidPipelineLogLines  = InvalidArgumentError("invalid pipeline log lines")
	InvalidPipelineLogID     = InvalidArgumentError("invalid pipeline log ID")
)
