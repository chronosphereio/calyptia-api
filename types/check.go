package types

// CheckStatus possible status of a check.
type CheckStatus string

const (
	CheckStatusOK      CheckStatus = "ok"
	CheckStatusFailed  CheckStatus = "failed"
	CheckStatusNew     CheckStatus = "new"
	CheckStatusRunning CheckStatus = "running"
)

// AllValidCheckStatuses all valid statuses for checks.
var AllValidCheckStatuses = [...]CheckStatus{
	CheckStatusNew,
	CheckStatusOK,
	CheckStatusRunning,
	CheckStatusFailed,
}

// AllValidCheckProtocols all valid protocols for checks.
var AllValidCheckProtocols = [...]PipelinePortProtocol{
	PipelineProtocolUDP,
	PipelineProtocolTCP,
}

// Check base check model.
type Check struct {
	Protocol PipelinePortProtocol `json:"protocol"`
	Status   CheckStatus          `json:"status"`
	Retries  uint                 `json:"retries"`
	Port     uint                 `json:"port"`
	Host     string               `json:"host"`
}

// ValidCheckProtocol checks if the given status is a valid check protocol.
func ValidCheckProtocol(p PipelinePortProtocol) bool {
	for _, protocol := range AllValidCheckProtocols {
		if p == protocol {
			return true
		}
	}
	return false
}

// ValidCheckStatus checks if the given string is a valid check status.
func ValidCheckStatus(s CheckStatus) bool {
	for _, status := range AllValidCheckStatuses {
		if s == status {
			return true
		}
	}
	return false
}
