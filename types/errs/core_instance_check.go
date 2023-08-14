package errs

const (
	InvalidCoreInstanceCheckStatus   = InvalidArgumentError("invalid core instance check status")
	InvalidCoreInstanceCheckID       = InvalidArgumentError("invalid core instance check ID")
	InvalidCoreInstanceCheckProtocol = InvalidArgumentError("invalid core instance check protocol")
	InvalidCoreInstanceCheckPort     = InvalidArgumentError("invalid core instance check port")
	InvalidCoreInstanceCheckHost     = InvalidArgumentError("invalid core instance check host")
	CoreInstanceCheckNotFound        = NotFoundError("core instance check not found")
)
