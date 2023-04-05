package types

type Health struct {
	Version  string         `json:"version"`
	Postgres HealthPostgres `json:"postgres"`
	Influx   HealthInflux   `json:"influx"`
}

type HealthPostgres struct {
	OK bool `json:"ok"`
}

type HealthInflux struct {
	OK bool `json:"ok"`
}
