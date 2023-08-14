package errs

const (
	InvalidStart    = InvalidArgumentError("invalid start")
	InvalidInterval = InvalidArgumentError("invalid interval")
	// InvalidMetricPlugin denotes that either `tag` nor `name` were not defined on the metric.
	InvalidMetricPlugin = InvalidArgumentError("invalid metric plugin")
	EmptyMetricTags     = InvalidArgumentError("empty metric tags")
	InvalidMetricTags   = InvalidArgumentError("invalid metric tags, missing valid tag property")
	MissingPipelineTag  = InvalidArgumentError("missing pipeline tag")
	ZeroPipelineIDs     = InvalidArgumentError("zero pipeline IDs")
)
