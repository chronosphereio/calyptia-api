package errs

const (
	InvalidStart    = InvalidArgumentError("invalid start")
	InvalidInterval = InvalidArgumentError("invalid interval")
	InvalidMetrics  = InvalidArgumentError("invalid metrics")
)
