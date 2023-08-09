package errs

const (
	InvalidPipelineFileID         = InvalidArgumentError("invalid pipeline file ID")
	InvalidPipelineFileName       = InvalidArgumentError("invalid pipeline file name")
	PipelineFilenameAlreadyExists = ConflictError("pipeline file name already exists")
	PipelineFileNotFound          = NotFoundError("pipeline file not found")
	PipelineFilesQuotaExceeded    = PermissionDeniedError("pipeline files quota exceeded")
	PipelineFileSizeExceeded      = InvalidArgumentError("pipeline file size exceeded")
)
