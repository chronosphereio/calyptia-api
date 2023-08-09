package errs

const (
	ClusterObjectAlreadyExists = InvalidArgumentError("cluster object already exists")
	InvalidClusterObjectID     = InvalidArgumentError("invalid cluster object ID")
	InvalidClusterObjectName   = InvalidArgumentError("invalid cluster object name")
	InvalidClusterObjectKind   = InvalidArgumentError("invalid cluster object kind")
	ClusterObjectNotFound      = NotFoundError("cluster object not found")
)
