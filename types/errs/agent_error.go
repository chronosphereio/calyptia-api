package errs

const (
	AgentErrorNotFound             = NotFoundError("agent error not found")
	InvalidAgentErrorID            = InvalidArgumentError("invalid agent error ID")
	InvalidAgentError              = InvalidArgumentError("invalid agent error")
	EitherAgentIDOrFleetIDRequired = InvalidArgumentError("either agent ID or fleet ID required")
	InvalidAgentErrorDismissReason = InvalidArgumentError("invalid agent error dismiss reason")
)
