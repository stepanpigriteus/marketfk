package model

type HealthStatus struct {
	Redis    string            `json:"redis"`
	Postgres string            `json:"postgres"`
	Ports    map[string]string `json:"ports"`
	Status   string            `json:"status"`
}
