package errs

const (
	InvalidClusterObjectRegexID     = InvalidArgumentError("invalid cluster object regex ID")
	InvalidClusterObjectRegex       = InvalidArgumentError("invalid cluster object regex")
	InvalidClusterObjectRegexDesc   = InvalidArgumentError("invalid cluster object regex description")
	ClusterObjectRegexAlreadyExists = ConflictError("cluster object regex already exists")
	ClusterObjectRegexNotFound      = NotFoundError("cluster object regex not found")
)
