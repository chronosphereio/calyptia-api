package errs

const (
	InvalidFleetFileID         = InvalidArgumentError("invalid fleet file ID")
	InvalidFleetFileName       = InvalidArgumentError("invalid fleet file name")
	FleetFilenameAlreadyExists = ConflictError("fleet file name already exists")
	FleetFileNotFound          = NotFoundError("fleet file not found")
	FleetFilesQuotaExceeded    = PermissionDeniedError("fleet files quota exceeded")
	FleetFileSizeExceeded      = InvalidArgumentError("fleet file size exceeded")
)
