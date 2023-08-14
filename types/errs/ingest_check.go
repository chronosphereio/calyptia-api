package errs

const (
	IngestCheckAlreadyExists  = InvalidArgumentError("ingestion check already exists")
	InvalidIngestCheckStatus  = InvalidArgumentError("invalid ingest check status")
	InvalidIngestCheckRetries = InvalidArgumentError("invalid ingest check retries")
	InvalidIngestCheckID      = InvalidArgumentError("invalid ingest check ID")
	InvalidConfigSectionKind  = InvalidArgumentError("invalid config section kind")
	IngestCheckNotFound       = NotFoundError("ingest check not found")
)
