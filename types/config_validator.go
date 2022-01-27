package types

// ValidatingConfig request body for validating a config.
type ValidatingConfig struct {
	Configs []ValidatingConfigEntry `json:"config"`
}

// ValidatingConfigEntry defines a single config to the validated. See `ValidatingConfig`.
type ValidatingConfigEntry struct {
	Command  string            `json:"command"`
	Name     string            `json:"name"`
	Optional map[string]string `json:"optional,omitempty"`
	ID       string            `json:"id"`
}

// ValidatedConfig response body after validating an agent config successfully.
type ValidatedConfig struct {
	Errors ConfigValidity `json:"errors"`
}

// ConfigValidity details.
type ConfigValidity struct {
	Runtime []string            `json:"runtime" `
	Input   map[string][]string `json:"input"`
	Output  map[string][]string `json:"output"`
	Filter  map[string][]string `json:"filter"`
}
