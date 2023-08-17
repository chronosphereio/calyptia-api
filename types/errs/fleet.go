package errs

const (
	InvalidFleetID           = InvalidArgumentError("invalid fleet ID")
	InvalidFleetName         = InvalidArgumentError("invalid fleet name")
	InvalidFluentBitVersion  = InvalidArgumentError("invalid fluentbit version")
	FleetNameTaken           = ConflictError("fleet name taken")
	FleetNotFound            = NotFoundError("fleet not found")
	InvalidFluentbitConfig   = InvalidArgumentError("invalid fluentbit config")
	InvalidFleetConfigFormat = InvalidArgumentError("invalid fleet config format")
)

func WrapFluentbitConfigError(err error) *DetailedError {
	if err == nil {
		return nil
	}

	msg := err.Error()
	return &DetailedError{Msg: "invalid fluentbit config", Detail: &msg, Parent: InvalidFluentbitConfig}
}
