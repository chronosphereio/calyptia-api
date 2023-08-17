package errs

const (
	InvalidPipelineKind                 = InvalidArgumentError("invalid pipeline kind")
	InvalidPipelineID                   = InvalidArgumentError("invalid pipeline ID")
	InvalidPipelineName                 = InvalidArgumentError("invalid pipeline name")
	InvalidPipelineMetadata             = InvalidArgumentError("invalid pipeline metadata")
	InvalidPipelineImage                = InvalidArgumentError("invalid pipeline image")
	PipelineNameAlreadyExists           = ConflictError("pipeline name already exists")
	CannotUpdatePipelineKind            = ConflictError("pipeline kind cannot be updated")
	PipelineNotFound                    = NotFoundError("pipeline not found")
	PipelinesQuotaExceeded              = PermissionDeniedError("pipelines quota exceeded")
	PipelineReplicasOnlyDeployments     = ConflictError("pipeline replicas can only be set for pipelines of kind deployment")
	InternalPipelinePermissionDenied    = PermissionDeniedError("pipeline can only be deleted by the system")
	ClusterLoggingPipelineAlreadyExists = ConflictError("a pipeline for cluster logging already exists")
	ZeroClusterObjectIDs                = InvalidArgumentError("zero cluster object IDs")
	ClusterLoggingOnlyNamespaces        = ConflictError("only cluster objects of type namespace can be associated to a cluster logging pipeline")
)