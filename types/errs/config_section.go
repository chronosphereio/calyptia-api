package errs

const (
	ConfigSectionNotFound        = NotFoundError("config section not found")
	InvalidConfigSectionID       = InvalidArgumentError("invalid config section ID")
	UnsupportedConfigSection     = InvalidArgumentError("unsupported config section")
	ConfigSectionSetSizeOverflow = InvalidArgumentError("config section set size overflow")
)
