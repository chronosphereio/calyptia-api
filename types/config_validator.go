package types

// ConfigValidity details.
type ConfigValidity struct {
	Runtime []string            `json:"runtime" `
	Input   map[string][]string `json:"input"`
	Output  map[string][]string `json:"output"`
	Filter  map[string][]string `json:"filter"`
}

// ValidatedConfig response body after validating an agent config successfully.
type ValidatedConfig struct {
	Errors ConfigValidity `json:"errors"`
}
