package types

type Health struct {
	Version  string         `json:"version"`
	Postgres HealthPostgres `json:"postgres"`
}

type HealthPostgres struct {
	OK bool `json:"ok"`
}
