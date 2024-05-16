package types

type ValidateFluentbitConfig struct {
	RawConfig    string       `json:"rawConfig"`
	ConfigFormat ConfigFormat `json:"configFormat"`
}
