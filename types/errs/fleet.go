package errs

const (
	InvalidFleetID           = InvalidArgumentError("invalid fleet ID")
	InvalidFleetName         = InvalidArgumentError("invalid fleet name")
	InvalidFluentBitVersion  = InvalidArgumentError("invalid fluentbit version")
	FleetNameTaken           = ConflictError("fleet name taken")
	FleetNotFound            = NotFoundError("fleet not found")
	InvalidFluentbitConfig   = InvalidArgumentError("invalid fluentbit config")
	InvalidFleetConfigFormat = InvalidArgumentError("invalid fleet config format")
)
