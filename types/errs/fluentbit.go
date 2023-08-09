package errs

const (
	PluginNotFound          = NotFoundError("plugin not found")
	InvalidConfigToValidate = InvalidArgumentError("invalid configuration validation payload")
)
