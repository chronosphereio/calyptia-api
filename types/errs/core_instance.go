package errs

const (
	InvalidCoreInstanceToken       = UnauthenticatedError("invalid aggregator token")
	InvalidAggregatorID            = InvalidArgumentError("invalid aggregator ID")
	InvalidCoreInstanceID          = InvalidArgumentError("invalid core instance ID")
	InvalidCoreInstanceName        = InvalidArgumentError("invalid aggregator name")
	InvalidCoreInstanceVersion     = InvalidArgumentError("invalid aggregator version")
	CoreInstanceGone               = GoneError("aggregator gone")
	CoreInstanceNotFound           = NotFoundError("aggregator not found")
	CoreInstanceNameAlreadyExists  = ConflictError("aggregator name already exists")
	CoreInstancesQuotaExceeded     = PermissionDeniedError("aggregators quota exceeded")
	InvalidCoreInstanceMetadata    = InvalidArgumentError("invalid aggregator metadata")
	InvalidCoreInstanceDockerImage = InvalidArgumentError("invalid docker image")
)
