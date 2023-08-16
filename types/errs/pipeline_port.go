package errs

import "fmt"

const (
	InvalidPipelinePortID       = InvalidArgumentError("invalid pipeline port ID")
	InvalidPipelinePortProtocol = InvalidArgumentError("invalid pipeline port protocol")
	InvalidPipelinePortNumber   = InvalidArgumentError("invalid pipeline port number")
	InvalidPipelinePortEndpoint = InvalidArgumentError("invalid pipeline port endpoint")
	PipelinePortNotFound        = NotFoundError("pipeline port not found")
	PipelinePortInUse           = ConflictError("pipeline port already in use")
)

func NewPipelinePortInUseError(protocol string, port uint) error {
	return NewDetailedError(PipelinePortInUse, fmt.Sprintf("pipeline port already in use: protocol %s, port %d", protocol, port))
}

func NewPipelinePortAlreadyInUseError(protocol string, port uint) *DetailedError {
	detail := fmt.Sprintf("protocol %s, port %d", protocol, port)
	return &DetailedError{Msg: "pipeline port already in use", Detail: &detail, Parent: Conflict}
}
