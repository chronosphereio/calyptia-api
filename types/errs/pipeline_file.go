package errs

const (
	InvalidPipelineFileID      = InvalidArgumentError("invalid pipeline file ID")
	PipelineFileAlreadyExists  = ConflictError("pipeline file already exists")
	PipelineFileNotFound       = NotFoundError("pipeline file not found")
	PipelineFilesQuotaExceeded = PermissionDeniedError("pipeline files quota exceeded")
)

const (
	InvalidFileName      = InvalidArgumentError("invalid file name")
	FileSizeExceeded     = InvalidArgumentError("file size exceeded")
	FileContentsRequired = InvalidArgumentError("file contents required")
)
