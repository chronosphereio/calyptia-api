package errs

const (
	InvalidTraceSessionID        = InvalidArgumentError("invalid trace session ID")
	TraceSessionNotFound         = NotFoundError("trace session not found")
	NoActiveTraceSession         = NotFoundError("no active trace session")
	ActiveTraceSessionInProgress = ConflictError("active trace session in progress")
	InvalidTraceLifespan         = InvalidArgumentError("invalid trace lifespan")
	ZeroTracePlugins             = InvalidArgumentError("zero trace plugins")
	InvalidTraceRecordKind       = InvalidArgumentError("invalid trace record kind")
	InvalidTraceRecords          = InvalidArgumentError("invalid trace record")
	TraceSessionTerminated       = GoneError("trace session terminated")
)
