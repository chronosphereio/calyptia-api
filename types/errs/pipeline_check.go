package errs

const (
	InvalidPipelineCheckStatus   = InvalidArgumentError("invalid pipeline check status")
	InvalidPipelineCheckID       = InvalidArgumentError("invalid pipeline check ID")
	InvalidPipelineCheckProtocol = InvalidArgumentError("invalid pipeline check protocol")
	InvalidPipelineCheckPort     = InvalidArgumentError("invalid pipeline check port")
	InvalidPipelineCheckHost     = InvalidArgumentError("invalid pipeline check host")
	PipelineCheckNotFound        = NotFoundError("pipeline check not found")
	PipelineCheckAlreadyExists   = ConflictError("pipeline check already exists")
)
