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

// ConfigValidityV2Property property details.
type ConfigValidityV2Property struct {
	ID       string   `json:"id"`
	Property string   `json:"property"`
	Text     string   `json:"text"`
	Errors   []string `json:"errors"`
}

// ConfigValidityV2Runtime runtime details.
type ConfigValidityV2Runtime struct {
	ID     string   `json:"id"`
	Errors []string `json:"errors"`
}

// ConfigValidityV2 details.
type ConfigValidityV2 struct {
	Runtime []ConfigValidityV2Runtime             `json:"runtime" `
	Input   map[string][]ConfigValidityV2Property `json:"input"`
	Output  map[string][]ConfigValidityV2Property `json:"output"`
	Filter  map[string][]ConfigValidityV2Property `json:"filter"`
}

// ValidatedConfigV2 response body after validating an agent config successfully against the v2 endpoint.
type ValidatedConfigV2 struct {
	Errors ConfigValidityV2 `json:"errors"`
}
